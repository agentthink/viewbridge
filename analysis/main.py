import os
import signal

# 解决 macOS 上 OpenMP 库冲突问题，必须在导入 cv2 等库之前设置
os.environ["KMP_DUPLICATE_LIB_OK"] = "TRUE"

import argparse
import base64
from concurrent import futures
import logging
import queue
import sys
import threading
import time
from typing import Any
import requests

import grpc

# 添加当前目录到 path 以支持直接运行
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

import logger
from detect import MotionDetector, ObjectDetector
from frame_capture import FrameCapture
import cv2

# 模型文件搜索候选路径（按优先级排序）
MODEL_SEARCH_PATHS = [
    ("../configs/owl.tflite", "tflite"),
    ("../configs/owl.onnx", "onnx"),
    ("./owl.tflite", "tflite"),
    ("./owl.onnx", "onnx"),
]

# 导入生成的 proto 代码
# 这些模块必须存在才能启动 gRPC 服务
import analysis_pb2
import analysis_pb2_grpc


slog = logging.getLogger("AI")

# 全局配置
GLOBAL_CONFIG = {
    "callback_url": "",
    "callback_secret": "",
}

# 保存父进程 PID，用于检测父进程是否退出
_PARENT_PID = os.getppid()


def _watch_parent_process():
    """
    监控父进程是否存活。当 Go 父进程退出后，Python 子进程应该自动退出，
    避免成为孤儿进程持续占用端口和资源。
    """
    while True:
        time.sleep(3)
        # 检查父进程是否还存在
        # 如果父进程退出，当前进程的 ppid 会变成 1 (init/launchd) 或其他进程
        current_ppid = os.getppid()
        if current_ppid != _PARENT_PID:
            slog.warning(
                f"父进程已退出 (原 PID: {_PARENT_PID}, 当前 PPID: {current_ppid})，Python 进程退出"
            )
            os._exit(0)


class CameraTask:
    def __init__(
        self,
        camera_id: str,
        rtsp_url: str,
        config: dict[str, Any],
        detector: ObjectDetector,
        motion_detector: MotionDetector,
    ) -> None:
        self.camera_id = camera_id
        self.rtsp_url = rtsp_url
        self.config = config
        self.detector = detector
        self.motion_detector = motion_detector

        self.status = "initializing"
        self.frames_processed = 0
        self.retry_count = 0
        self.last_error = ""
        self._stop_event = threading.Event()
        self._thread: threading.Thread | None = None

        # 告警去重缓存：key=label, value=(last_alert_time, last_box)
        # 冷却期内同一目标（标签相同 + IoU > 0.3）不重复告警
        self._alert_cache: dict[str, tuple[float, dict]] = {}
        self._cooldown: float = config.get("alert_cooldown_seconds", 30.0)

        self.frame_queue = queue.Queue(maxsize=1)
        self.capture = FrameCapture(
            rtsp_url,
            self.frame_queue,
            config.get("detect_interval_seconds", 5.0),
            config.get("retry_limit", 10),
        )

    def start(self):
        self.status = "running"
        self.capture.start()
        self._stop_event.clear()
        self._thread = threading.Thread(target=self._analysis_loop, daemon=True)
        self._thread.start()
        slog.info(f"CameraTask started for {self.camera_id}")

    def stop(self):
        self.status = "stopping"
        self._stop_event.set()
        self.capture.stop()
        if self._thread:
            self._thread.join(timeout=2)
        slog.info(f"CameraTask stopped for {self.camera_id}")

    def _analysis_loop(self):
        error_streak = 0
        retry_limit = int(self.config.get("retry_limit", 10))

        while not self._stop_event.is_set():
            # 检查 FrameCapture 是否已达到重试上限
            if self.capture.is_failed:
                self.status = "error"
                self.last_error = self.capture.last_error
                self._send_stopped_callback("capture_failed", self.last_error)
                slog.error(
                    f"CameraTask {self.camera_id} 因帧捕获失败而停止: {self.last_error}"
                )
                break

            try:
                try:
                    frame = self.frame_queue.get(timeout=2.0)
                except queue.Empty:
                    slog.debug("CameraTask frame queue empty, skipping")
                    continue

                error_streak = 0
                self.frames_processed += 1

                zones = self.config.get("zones") or []
                motion_boxes, has_motion = self.motion_detector.detect(
                    frame, self.camera_id, zones if zones else None
                )

                if not has_motion:
                    continue

                try:
                    detections = self._detect_with_zones(frame, zones)
                except Exception as e:
                    slog.error(f"CameraTask detect error: {e}")
                    continue

                if not detections:
                    continue

                detections = self._dedup(detections)
                if not detections:
                    continue
                self._send_detection_callback(detections, frame)
            except Exception as e:
                slog.error(f"CameraTask analysis loop error: {e}")
                error_streak += 1
                self.last_error = str(e)
                if error_streak >= retry_limit:
                    self.status = "error"
                    self._send_stopped_callback("error", self.last_error)
                    self.capture.stop()
                    break
                # 防止 cpu 在异常里空转
                time.sleep(1)

    def _detect_with_zones(self, frame, zones: list[dict]) -> list[dict]:
        """
        多区域目标检测：有区域配置时逐区域过滤，无区域时全图检测。
        每个区域使用自己的 labels；检测框中心点必须在多边形内部才计入。
        """
        import numpy as np
        import cv2

        threshold = self.config.get("threshold", 0.5)
        global_labels = self.config.get("labels") or []
        safe_global = [str(l) for l in global_labels] if global_labels else None

        if not zones:
            dets, _ = self.detector.detect(frame, threshold=threshold, label_filter=safe_global)
            return dets

        h, w = frame.shape[:2]
        result = []
        seen_boxes: set[tuple] = set()

        for zone in zones:
            pts_flat = zone.get("points", [])
            zone_labels = zone.get("labels") or global_labels
            safe_labels = [str(l) for l in zone_labels] if zone_labels else safe_global

            # 将归一化坐标转为像素多边形
            poly = None
            if pts_flat and len(pts_flat) >= 6:
                pts = [(int(pts_flat[i] * w), int(pts_flat[i + 1] * h)) for i in range(0, len(pts_flat), 2)]
                poly = np.array(pts, dtype=np.float32)

            dets, _ = self.detector.detect(frame, threshold=threshold, label_filter=safe_labels)

            for det in dets:
                box = det["box"]
                # 判断：检测框底部中心点或框中心点，任意一个在多边形内则计入
                # 底部中心点适合人物检测（人的脚在区域内即视为进入区域）
                if poly is not None:
                    cx = (box["x_min"] + box["x_max"]) / 2.0
                    cy_center = (box["y_min"] + box["y_max"]) / 2.0
                    cy_bottom = float(box["y_max"])
                    in_zone = (
                        cv2.pointPolygonTest(poly, (cx, cy_bottom), measureDist=False) >= 0
                        or cv2.pointPolygonTest(poly, (cx, cy_center), measureDist=False) >= 0
                    )
                    if not in_zone:
                        continue

                key = (det["label"], box["x_min"], box["y_min"], box["x_max"], box["y_max"])
                if key not in seen_boxes:
                    seen_boxes.add(key)
                    result.append(det)

        return result

    @staticmethod
    def _iou(a: dict, b: dict) -> float:
        """
        计算两个检测框的 IoU（交并比）。
        用于判断两次检测是否为同一目标，避免因目标轻微移动导致重复告警。
        """
        ax1, ay1, ax2, ay2 = a["x_min"], a["y_min"], a["x_max"], a["y_max"]
        bx1, by1, bx2, by2 = b["x_min"], b["y_min"], b["x_max"], b["y_max"]

        inter_x1 = max(ax1, bx1)
        inter_y1 = max(ay1, by1)
        inter_x2 = min(ax2, bx2)
        inter_y2 = min(ay2, by2)

        inter_area = max(0, inter_x2 - inter_x1) * max(0, inter_y2 - inter_y1)
        if inter_area == 0:
            return 0.0

        area_a = (ax2 - ax1) * (ay2 - ay1)
        area_b = (bx2 - bx1) * (by2 - by1)
        union_area = area_a + area_b - inter_area
        return inter_area / union_area if union_area > 0 else 0.0

    def _dedup(self, detections: list[dict]) -> list[dict]:
        """
        告警去重（方案 B：时间窗口 + IoU 空间比对）。
        冷却期内，标签相同且 IoU > 0.3 的目标视为重复，过滤掉不告警。
        冷却期外，或空间位置差异大（新目标），放行并刷新缓存。
        """
        now = time.time()
        result = []
        for det in detections:
            label = det["label"]
            box = det["box"]
            cached = self._alert_cache.get(label)
            if cached is not None:
                last_time, last_box = cached
                if now - last_time < self._cooldown and self._iou(box, last_box) > 0.3:
                    slog.debug(
                        f"[dedup] 抑制重复告警: camera={self.camera_id} label={label} "
                        f"cooldown_remaining={self._cooldown - (now - last_time):.1f}s"
                    )
                    continue
            self._alert_cache[label] = (now, box)
            result.append(det)
        return result

    def _send_detection_callback(self, detections, frame):
        timestamp = int(time.time() * 1000)
        draw_frame = frame.copy()
        for det in detections:
            box = det["box"]
            label = f"{det['label']} {det['confidence']:.2f}"

            # 坐标
            p1 = (box["x_min"], box["y_min"])
            p2 = (box["x_max"], box["y_max"])

            # 画矩形框 (红色，线宽2)
            cv2.rectangle(draw_frame, p1, p2, (0, 0, 255), 2)

            # 画文字背景条，防止文字看不清
            t_size = cv2.getTextSize(label, cv2.FONT_HERSHEY_SIMPLEX, 0.5, 1)[0]
            p2_text = (p1[0] + t_size[0], p1[1] - t_size[1] - 3)
            cv2.rectangle(draw_frame, p1, p2_text, (0, 0, 255), -1)  # -1 表示实心填充

            # 画文字 (白色)
            cv2.putText(
                draw_frame,
                label,
                (p1[0], p1[1] - 2),
                cv2.FONT_HERSHEY_SIMPLEX,
                0.5,
                (255, 255, 255),
                1,
            )
        success, buffer = cv2.imencode(".jpg", draw_frame)
        snapshot_b64 = ""
        if success:
            snapshot_b64 = base64.b64encode(buffer).decode("utf-8")

        payload = {
            "camera_id": self.camera_id,
            "timestamp": timestamp,
            "detections": detections,
            "snapshot": snapshot_b64,
            "snapshot_width": frame.shape[1],
            "snapshot_height": frame.shape[0],
        }

        send_callback(self.config, "/events", payload)

    def _send_stopped_callback(self, reason, message):
        payload = {
            "camera_id": self.camera_id,
            "timestamp": int(time.time() * 1000),
            "reason": reason,
            "message": message,
        }
        send_callback(self.config, "/stopped", payload)


class HealthServicer(analysis_pb2_grpc.HealthServicer):
    def __init__(self, servicer):
        self._servicer = servicer

    def Check(self, request, context):
        if not self._servicer.is_ready:
            return analysis_pb2.HealthCheckResponse(
                status=analysis_pb2.HealthCheckResponse.NOT_SERVING
            )
        return analysis_pb2.HealthCheckResponse(
            status=analysis_pb2.HealthCheckResponse.SERVING
        )


class AnalysisServiceServicer(analysis_pb2_grpc.AnalysisServiceServicer):
    def __init__(self, model_path):
        self._camera_tasks: dict[str, CameraTask] = {}
        self._lock = threading.Lock()
        self._is_ready = False
        self._start_time = time.time()

        self.object_detector = ObjectDetector(model_path)
        self.motion_detector = MotionDetector()

    def is_ready(self) -> bool:
        return self._is_ready

    def initialize(self):
        slog.info("AnalysisService initializing...")
        success = self.object_detector.load_model()
        self._is_ready = success

        if not success:
            slog.error("AnalysisService initialization failed")
            return
        slog.info("AnalysisService initialized")
        threading.Thread(target=send_started_callback).start()

    def StartCamera(self, request, context):
        if not self._is_ready:
            context.set_details("model loadding")
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            return analysis_pb2.StartCameraResponse(
                success=False, message="model loadding"
            )
        camera_id = request.camera_id
        with self._lock:
            if camera_id in self._camera_tasks:
                slog.info(
                    f"Camera {camera_id} already exists, status: {self._camera_tasks[camera_id].status}"
                )
                return analysis_pb2.StartCameraResponse(
                    success=True, message="任务已运行"
                )
            cb_url = request.callback_url or GLOBAL_CONFIG["callback_url"]
            cb_secret = request.callback_secret or GLOBAL_CONFIG["callback_secret"]
            if not cb_url:
                context.set_details("callback url is required")
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                return analysis_pb2.StartCameraResponse(
                    success=False, message="callback url is required"
                )
            # 将 proto AnalysisZone 列表转为纯 Python dict 存入 config
            zones = [
                {"points": list(z.points), "labels": list(z.labels), "name": z.name}
                for z in request.zones
            ]
            config = {
                "detect_interval_seconds": request.detect_interval_seconds,
                "alert_cooldown_seconds": request.alert_cooldown_seconds,
                "labels": list(request.labels),
                "threshold": request.threshold,
                "zones": zones,
                "retry_limit": request.retry_limit,
                "callback_url": cb_url,
                "callback_secret": cb_secret,
            }

            task = CameraTask(
                camera_id,
                rtsp_url=request.rtsp_url,
                config=config,
                detector=self.object_detector,
                motion_detector=self.motion_detector,
            )
            task.start()
            self._camera_tasks[camera_id] = task

        timeout = 5.0
        start = time.time()
        w, h, fps = 0, 0, 0.0
        while time.time() - start < timeout:
            w, h, fps = task.capture.get_stream_info()
            if w > 0:
                break
            time.sleep(0.5)
        return analysis_pb2.StartCameraResponse(
            success=True,
            message="任务已启动",
            source_width=w,
            source_height=h,
            source_fps=fps,
        )

    def StopCamera(self, request, context):
        camera_id = request.camera_id
        with self._lock:
            if camera_id not in self._camera_tasks:
                return analysis_pb2.StopCameraResponse(
                    success=False, message="Camera not found"
                )

            task = self._camera_tasks.pop(camera_id)
            task.stop()
            return analysis_pb2.StopCameraResponse(success=True, message="任务已停止")

    def GetStatus(self, request, context):
        response = analysis_pb2.StatusResponse()
        response.is_ready = self._is_ready
        response.stats.active_streams = len(self._camera_tasks)
        response.stats.uptime_seconds = int(time.time() - self._start_time)

        with self._lock:
            for cid, task in self._camera_tasks.items():
                cam_status = analysis_pb2.CameraStatus(
                    camera_id=cid,
                    status=task.status,
                    frames_processed=task.frames_processed,
                    retry_count=task.retry_count,
                    last_error=task.last_error,
                )
                response.cameras.append(cam_status)
        return response


def send_callback(config: dict, path: str, payload: dict):
    """
    发送回调到指定路径，路径会拼接到 callback_url 后面。
    例如: callback_url=http://127.0.0.1:15123/webhook, path=/events
    最终请求: POST http://127.0.0.1:15123/webhook/events
    secret 通过 Secret header 传递，不出现在 URL 中
    """
    url = config.get("callback_url", "")
    secret = config.get("callback_secret", "")
    if not url:
        return

    full_url = url.rstrip("/") + path
    headers = {"Content-Type": "application/json"}
    if secret:
        headers["Secret"] = secret

    try:
        threading.Thread(
            target=requests.post,
            args=(full_url,),
            kwargs={
                "json": payload,
                "headers": headers,
                "timeout": 5.0,
            },
        ).start()
    except Exception as e:
        slog.error(f"Failed to send callback to {path}: {e}")


def send_started_callback():
    """
    向 Go 服务发送启动通知，用于确认 Python 进程与 Go 服务的连接是否正常。
    如果 Go 服务返回 404，说明回调接口不存在，Python 进程应该退出，避免成为孤儿进程。
    """
    url = GLOBAL_CONFIG["callback_url"]
    secret = GLOBAL_CONFIG["callback_secret"]
    if not url:
        return

    full_url = url.rstrip("/") + "/started"
    headers = {"Content-Type": "application/json"}
    if secret:
        headers["Secret"] = secret

    payload = {
        "timestamp": int(time.time() * 1000),
        "message": "AI Analysis Service Started",
    }

    max_retries = 3
    retry_interval = 2

    for attempt in range(1, max_retries + 1):
        slog.info(f"Sending started callback (attempt {attempt}/{max_retries})...")
        try:
            resp = requests.post(full_url, json=payload, headers=headers, timeout=5)
            if resp.status_code == 404 and attempt == max_retries - 1:
                slog.error(f"回调接口返回 404，Go 服务可能已停止，退出 Python 进程")
                os._exit(1)
            if resp.ok:
                slog.info("启动通知发送成功")
                return
            slog.warning(f"启动通知返回非成功状态: {resp.status_code} {full_url}")
        except requests.exceptions.ConnectionError as e:
            slog.warning(f"发送启动通知失败 (连接错误): {e}")
        except Exception as e:
            slog.error(f"发送启动通知失败: {e}")

        if attempt < max_retries:
            time.sleep(retry_interval)

    slog.error(f"启动通知发送失败，已重试 {max_retries} 次")


def send_keepalive_callback(stats: dict):
    """
    发送心跳回调，用于定期向 Go 服务报告 AI 服务状态。
    """
    url = GLOBAL_CONFIG["callback_url"]
    secret = GLOBAL_CONFIG["callback_secret"]
    if not url:
        return

    full_url = url.rstrip("/") + "/keepalive"
    headers = {"Content-Type": "application/json"}
    if secret:
        headers["Secret"] = secret

    payload = {
        "timestamp": int(time.time() * 1000),
        "stats": stats,
        "message": "Service running normally",
    }

    try:
        requests.post(full_url, json=payload, headers=headers, timeout=5)
    except Exception as e:
        slog.debug(f"Failed to send keepalive callback: {e}")


def serve(port, model_path):
    # 启动父进程监控线程，确保 Go 退出时 Python 也退出
    threading.Thread(target=_watch_parent_process, daemon=True).start()

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=20))
    servicer = AnalysisServiceServicer(model_path)
    analysis_pb2_grpc.add_AnalysisServiceServicer_to_server(servicer, server)

    health_servicer = HealthServicer(servicer)
    analysis_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)

    server.add_insecure_port(f"[::]:{port}")
    server.start()
    slog.info(f"AnalysisService started: 0.0.0.0:{port}")

    threading.Thread(target=servicer.initialize).start()

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop(0)


def discover_model(model_arg: str) -> str:
    """
    自动发现可用模型文件
    优先级：../configs/owl.tflite > ../configs/owl.onnx > ./owl.tflite > ./owl.onnx > 命令行参数
    """
    script_dir = os.path.dirname(os.path.abspath(__file__))

    for rel_path, _ in MODEL_SEARCH_PATHS:
        full_path = os.path.normpath(os.path.join(script_dir, rel_path))
        if os.path.exists(full_path):
            slog.info(f"发现模型文件: {full_path}")
            return full_path

    # 回退到命令行参数指定的模型
    if os.path.isabs(model_arg):
        return model_arg

    # 相对路径基于脚本目录解析
    return os.path.normpath(os.path.join(script_dir, model_arg))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--port", type=int, default=50051)
    parser.add_argument("--model", type=str, default="owl.onnx")
    parser.add_argument(
        "--callback-url",
        type=str,
        default="http://127.0.0.1:15123",
        help="回调基础URL，各回调路由会自动拼接",
    )
    parser.add_argument("--callback-secret", type=str, default="", help="回调秘钥")
    parser.add_argument(
        "--log-level",
        type=str,
        default="INFO",
        help="日志级别 (DEBUG/INFO/ERROR)",
    )
    args = parser.parse_args()
    logger.setup_logging(level_str=args.log_level)

    GLOBAL_CONFIG["callback_url"] = args.callback_url
    GLOBAL_CONFIG["callback_secret"] = args.callback_secret

    # 自动发现模型文件
    model_path = discover_model(args.model)

    slog.debug(
        f"log level: {args.log_level}, model: {model_path}, callback url: {args.callback_url}, callback secret: {args.callback_secret}"
    )

    serve(args.port, model_path)


if __name__ == "__main__":
    main()
