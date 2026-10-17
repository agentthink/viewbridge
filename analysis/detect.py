import logging
import os
import time
from abc import ABC, abstractmethod
from typing import Any

import numpy as np
import cv2

slog = logging.getLogger("Detector")

# COCO 数据集 80 类标签
COCO_LABELS = [
    "person",
    "bicycle",
    "car",
    "motorcycle",
    "airplane",
    "bus",
    "train",
    "truck",
    "boat",
    "traffic light",
    "fire hydrant",
    "stop sign",
    "parking meter",
    "bench",
    "bird",
    "cat",
    "dog",
    "horse",
    "sheep",
    "cow",
    "elephant",
    "bear",
    "zebra",
    "giraffe",
    "backpack",
    "umbrella",
    "handbag",
    "tie",
    "suitcase",
    "frisbee",
    "skis",
    "snowboard",
    "sports ball",
    "kite",
    "baseball bat",
    "baseball glove",
    "skateboard",
    "surfboard",
    "tennis racket",
    "bottle",
    "wine glass",
    "cup",
    "fork",
    "knife",
    "spoon",
    "bowl",
    "banana",
    "apple",
    "sandwich",
    "orange",
    "broccoli",
    "carrot",
    "hot dog",
    "pizza",
    "donut",
    "cake",
    "chair",
    "couch",
    "potted plant",
    "bed",
    "dining table",
    "toilet",
    "tv",
    "laptop",
    "mouse",
    "remote",
    "keyboard",
    "cell phone",
    "microwave",
    "oven",
    "toaster",
    "sink",
    "refrigerator",
    "book",
    "clock",
    "vase",
    "scissors",
    "teddy bear",
    "hair drier",
    "toothbrush",
]


class ModelBackend(ABC):
    """
    模型推理后端抽象接口
    不同模型格式（ONNX、TFLite）需实现此接口，确保上层调用逻辑统一
    """

    @abstractmethod
    def load(self, model_path: str) -> bool:
        """加载模型文件，返回是否成功"""
        pass

    @abstractmethod
    def is_ready(self) -> bool:
        """检查模型是否已加载并可用"""
        pass

    @abstractmethod
    def get_input_shape(self) -> tuple:
        """获取模型输入形状，用于预处理"""
        pass

    @abstractmethod
    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        """执行推理，返回原始输出"""
        pass


class ONNXBackend(ModelBackend):
    """
    ONNX Runtime 推理后端
    使用 onnxruntime 库加载和执行 ONNX 格式模型
    """

    def __init__(self):
        self.session = None
        self.input_name: str = ""
        self.input_shape: tuple = (1, 3, 640, 640)
        self._is_ready = False

    def load(self, model_path: str) -> bool:
        try:
            import onnxruntime as ort

            slog.info(f"加载 ONNX 模型: {model_path} ...")
            start_time = time.time()

            sess_options = ort.SessionOptions()
            sess_options.graph_optimization_level = (
                ort.GraphOptimizationLevel.ORT_ENABLE_ALL
            )
            sess_options.intra_op_num_threads = 4
            sess_options.inter_op_num_threads = 2

            providers = ["CPUExecutionProvider"]
            self.session = ort.InferenceSession(
                model_path, sess_options=sess_options, providers=providers
            )

            input_info = self.session.get_inputs()[0]
            self.input_name = input_info.name
            self.input_shape = tuple(input_info.shape)

            elapsed = time.time() - start_time
            slog.info(
                f"ONNX 模型加载完成 (耗时: {elapsed:.2f}s, 输入形状: {self.input_shape})"
            )
            self._is_ready = True
            return True
        except ImportError:
            slog.error("未安装 onnxruntime，无法加载 ONNX 模型")
            return False
        except Exception as e:
            slog.error(f"加载 ONNX 模型失败: {e}")
            return False

    def is_ready(self) -> bool:
        return self._is_ready and self.session is not None

    def get_input_shape(self) -> tuple:
        return self.input_shape

    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        if not self.session:
            raise RuntimeError("ONNX 模型未加载")
        outputs = self.session.run(None, {self.input_name: input_tensor})
        return np.asarray(outputs[0])


class TFLiteBackend(ModelBackend):
    """
    TensorFlow Lite 推理后端
    支持两种模型格式：
    1. YOLO 格式：单输出张量 (1, 84, 8400)
    2. SSD 格式：多输出张量（boxes, classes, scores, num_detections）
    """

    def __init__(self):
        self.interpreter: Any = None
        self.input_details: list[dict[str, Any]] = []
        self.output_details: list[dict[str, Any]] = []
        self.input_shape: tuple = (1, 640, 640, 3)
        self._is_ready = False
        self._is_nhwc = True
        self._is_ssd_format = False  # 区分 SSD 和 YOLO 格式
        self._input_quantization: tuple[float, int] = (1.0, 0)  # scale, zero_point

    def load(self, model_path: str) -> bool:
        try:
            Interpreter = None
            try:
                from tflite_runtime.interpreter import Interpreter  # type: ignore
            except ImportError:
                try:
                    from ai_edge_litert.interpreter import Interpreter  # type: ignore
                except ImportError:
                    try:
                        import tensorflow as tf

                        Interpreter = tf.lite.Interpreter
                    except ImportError:
                        pass

            if Interpreter is None:
                raise ImportError("未找到 tflite_runtime、ai_edge_litert 或 tensorflow")

            slog.info(f"加载 TFLite 模型: {model_path} ...")
            start_time = time.time()

            self.interpreter = Interpreter(model_path=model_path)
            self.interpreter.allocate_tensors()

            self.input_details = self.interpreter.get_input_details()
            self.output_details = self.interpreter.get_output_details()

            input_shape = self.input_details[0]["shape"]
            self.input_shape = tuple(input_shape)

            if len(self.input_shape) == 4:
                self._is_nhwc = self.input_shape[3] == 3

            # 获取输入量化参数（用于 uint8 量化模型）
            quant_params = self.input_details[0].get("quantization_parameters", {})
            scales = quant_params.get("scales", np.array([1.0]))
            zero_points = quant_params.get("zero_points", np.array([0]))
            if len(scales) > 0 and len(zero_points) > 0:
                self._input_quantization = (float(scales[0]), int(zero_points[0]))

            # 检测模型格式：SSD 模型通常有4个输出（boxes, classes, scores, num）
            # 且输出名称包含 "TFLite_Detection_PostProcess"
            self._is_ssd_format = len(self.output_details) >= 3 and any(
                "Detection" in d.get("name", "") for d in self.output_details
            )

            elapsed = time.time() - start_time
            format_name = "SSD" if self._is_ssd_format else "YOLO"
            slog.info(
                f"TFLite 模型加载完成 (耗时: {elapsed:.2f}s, 输入: {self.input_shape}, "
                f"格式: {format_name}, 量化: {self._input_quantization})"
            )
            self._is_ready = True
            return True
        except ImportError:
            slog.error("未安装 tflite_runtime 或 tensorflow，无法加载 TFLite 模型")
            return False
        except Exception as e:
            slog.error(f"加载 TFLite 模型失败: {e}")
            return False

    def is_ready(self) -> bool:
        return self._is_ready and self.interpreter is not None

    def get_input_shape(self) -> tuple:
        return self.input_shape

    def is_nhwc(self) -> bool:
        """返回模型是否使用 NHWC 格式"""
        return self._is_nhwc

    def is_ssd_format(self) -> bool:
        """返回是否为 SSD 格式（多输出张量）"""
        return self._is_ssd_format

    def get_input_quantization(self) -> tuple[float, int]:
        """返回输入量化参数 (scale, zero_point)"""
        return self._input_quantization

    def get_input_dtype(self) -> np.dtype:
        """返回模型期望的输入数据类型"""
        return self.input_details[0]["dtype"]

    def infer(self, input_tensor: np.ndarray) -> np.ndarray:
        """执行推理，返回第一个输出张量（用于 YOLO 格式）"""
        if not self.interpreter or len(self.input_details) == 0:
            raise RuntimeError("TFLite 模型未加载")

        input_dtype = self.input_details[0]["dtype"]
        if input_tensor.dtype != input_dtype:
            input_tensor = input_tensor.astype(input_dtype)

        self.interpreter.set_tensor(self.input_details[0]["index"], input_tensor)
        self.interpreter.invoke()

        output = self.interpreter.get_tensor(self.output_details[0]["index"])
        return np.asarray(output)

    def infer_ssd(
        self, input_tensor: np.ndarray
    ) -> tuple[np.ndarray, np.ndarray, np.ndarray, int]:
        """
        执行 SSD 格式推理，返回解析后的检测结果
        SSD 输出格式（已内置后处理）：
        - boxes: (1, num_boxes, 4) 归一化坐标 [y_min, x_min, y_max, x_max]
        - classes: (1, num_boxes) 类别 ID
        - scores: (1, num_boxes) 置信度分数
        - num_detections: 有效检测数量
        """
        if not self.interpreter or len(self.input_details) == 0:
            raise RuntimeError("TFLite 模型未加载")

        input_dtype = self.input_details[0]["dtype"]
        if input_tensor.dtype != input_dtype:
            input_tensor = input_tensor.astype(input_dtype)

        self.interpreter.set_tensor(self.input_details[0]["index"], input_tensor)
        self.interpreter.invoke()

        # 按名称或索引获取各输出张量
        boxes = None
        classes = None
        scores = None
        num_detections = 0

        for detail in self.output_details:
            name = detail.get("name", "")
            tensor = self.interpreter.get_tensor(detail["index"])

            if "boxes" in name.lower() or (
                detail["shape"][-1] == 4 and len(detail["shape"]) == 3
            ):
                boxes = np.asarray(tensor)
            elif "class" in name.lower() or (
                len(detail["shape"]) == 2
                and detail["shape"][1] > 1
                and boxes is not None
            ):
                classes = np.asarray(tensor)
            elif "score" in name.lower() or ":2" in name:
                scores = np.asarray(tensor)
            elif "num" in name.lower() or (
                len(detail["shape"]) == 1 and detail["shape"][0] == 1
            ):
                num_detections = int(tensor[0])

        # 兜底处理：按输出顺序分配
        if boxes is None or classes is None or scores is None:
            outputs = [
                self.interpreter.get_tensor(d["index"]) for d in self.output_details
            ]
            if len(outputs) >= 4:
                boxes = np.asarray(outputs[0])
                classes = np.asarray(outputs[1])
                scores = np.asarray(outputs[2])
                num_detections = int(outputs[3][0]) if outputs[3].size > 0 else 0

        if boxes is None:
            boxes = np.array([])
        if classes is None:
            classes = np.array([])
        if scores is None:
            scores = np.array([])

        return boxes, classes, scores, num_detections


def get_model_type(model_path: str) -> str:
    """根据模型文件后缀判断模型类型"""
    ext = os.path.splitext(model_path)[1].lower()
    return "tflite" if ext == ".tflite" else "onnx"


def create_backend(model_type: str) -> ModelBackend:
    """
    根据模型类型创建对应的推理后端
    """
    if model_type == "tflite":
        return TFLiteBackend()
    else:
        return ONNXBackend()


class ObjectDetector:
    """
    目标检测器 - 支持多种模型格式（ONNX、TFLite）
    通过统一的 ModelBackend 接口实现模型无关的检测逻辑
    """

    def __init__(self, model_path: str):
        self.model_path = model_path
        self.model_type = get_model_type(model_path)
        self.backend: ModelBackend | None = None
        self.input_shape: tuple = (1, 3, 640, 640)
        self._is_ready = False
        self.names: dict[int, str] = {i: name for i, name in enumerate(COCO_LABELS)}

    def load_model(self) -> bool:
        """加载模型并初始化推理后端"""
        try:
            start_time = time.time()

            # 创建对应类型的后端
            self.backend = create_backend(self.model_type)

            # 加载模型
            if not self.backend.load(self.model_path):
                return False

            self.input_shape = self.backend.get_input_shape()

            # 预热模型
            self._warmup()

            elapsed = time.time() - start_time
            slog.info(f"模型预热完成 (总耗时: {elapsed:.2f}s)")
            self._is_ready = True
            return True
        except Exception as e:
            slog.error(f"加载模型失败: {e}")
            return False

    def _warmup(self) -> None:
        """预热模型，减少首次推理延迟"""
        if not self.backend:
            return

        dummy_img = np.zeros((640, 640, 3), dtype=np.uint8)
        dummy_input = self._preprocess(dummy_img)
        self.backend.infer(dummy_input)
        slog.info("模型预热完成")

    def is_ready(self) -> bool:
        return self._is_ready and self.backend is not None and self.backend.is_ready()

    def _get_target_size(self) -> int:
        """获取模型期望的输入尺寸"""
        # NCHW: (1, 3, H, W) -> shape[2]
        # NHWC: (1, H, W, 3) -> shape[1]
        if self._is_nhwc_format():
            return int(self.input_shape[1])
        return int(self.input_shape[2])

    def _is_nhwc_format(self) -> bool:
        """判断当前后端是否使用 NHWC 格式"""
        if isinstance(self.backend, TFLiteBackend):
            return self.backend.is_nhwc()
        return False

    def _is_ssd_format(self) -> bool:
        """判断当前后端是否为 SSD 格式"""
        if isinstance(self.backend, TFLiteBackend):
            return self.backend.is_ssd_format()
        return False

    def _preprocess(self, image: np.ndarray) -> np.ndarray:
        """
        预处理图像：调整大小、归一化、转换格式
        根据后端类型自动选择 NCHW 或 NHWC 格式，并处理量化输入
        """
        target_size = self._get_target_size()
        h, w = image.shape[:2]

        # SSD 模型使用直接缩放（不保持宽高比）
        # YOLO 模型使用 letterbox 缩放（保持宽高比）
        if self._is_ssd_format():
            resized = cv2.resize(
                image, (target_size, target_size), interpolation=cv2.INTER_LINEAR
            )
            rgb = resized[:, :, ::-1]  # BGR -> RGB

            # 检查是否需要量化为 uint8
            if isinstance(self.backend, TFLiteBackend):
                input_dtype = self.backend.get_input_dtype()
                if input_dtype == np.uint8:
                    # 直接返回 uint8 格式
                    return np.expand_dims(rgb.astype(np.uint8), axis=0)

            # float32 格式
            rgb = rgb.astype(np.float32) / 255.0
            return np.expand_dims(rgb, axis=0)

        # YOLO letterbox 预处理
        scale = min(target_size / h, target_size / w)
        new_h, new_w = int(h * scale), int(w * scale)

        resized = cv2.resize(image, (new_w, new_h), interpolation=cv2.INTER_LINEAR)

        canvas = np.full((target_size, target_size, 3), 114, dtype=np.uint8)
        top = (target_size - new_h) // 2
        left = (target_size - new_w) // 2
        canvas[top : top + new_h, left : left + new_w] = resized

        rgb = canvas[:, :, ::-1].astype(np.float32) / 255.0

        if self._is_nhwc_format():
            return np.expand_dims(rgb, axis=0)
        else:
            blob = rgb.transpose(2, 0, 1)
            return np.expand_dims(blob, axis=0)

    def _postprocess(
        self,
        outputs: np.ndarray,
        original_shape: tuple[int, int],
        threshold: float,
        label_filter: list[str] | None = None,
    ) -> list[dict[str, Any]]:
        """
        后处理 YOLO 输出：解析检测框、应用 NMS、坐标转换
        YOLO11 输出格式: (1, 84, 8400) -> 84 = 4 (bbox) + 80 (classes)
        """
        # 转置为 (8400, 84) 便于处理
        predictions = outputs[0].T  # (8400, 84)

        # 提取边界框和类别分数
        boxes = predictions[:, :4]  # x_center, y_center, width, height
        scores = predictions[:, 4:]  # 80 个类别的分数

        # 获取每个检测框的最高分数和对应类别
        class_ids = np.argmax(scores, axis=1)
        confidences = np.max(scores, axis=1)

        # 过滤低置信度检测
        mask = confidences >= threshold
        boxes = boxes[mask]
        confidences = confidences[mask]
        class_ids = class_ids[mask]

        if len(boxes) == 0:
            return []

        # 转换坐标：center_x, center_y, w, h -> x1, y1, x2, y2
        x_center, y_center, w, h = boxes[:, 0], boxes[:, 1], boxes[:, 2], boxes[:, 3]
        x1 = x_center - w / 2
        y1 = y_center - h / 2
        x2 = x_center + w / 2
        y2 = y_center + h / 2

        # 缩放坐标到原始图像尺寸
        orig_h, orig_w = original_shape
        target_size = self._get_target_size()
        scale = min(target_size / orig_h, target_size / orig_w)
        pad_h = (target_size - orig_h * scale) / 2
        pad_w = (target_size - orig_w * scale) / 2

        x1 = (x1 - pad_w) / scale
        y1 = (y1 - pad_h) / scale
        x2 = (x2 - pad_w) / scale
        y2 = (y2 - pad_h) / scale

        # 裁剪到图像边界
        x1 = np.clip(x1, 0, orig_w)
        y1 = np.clip(y1, 0, orig_h)
        x2 = np.clip(x2, 0, orig_w)
        y2 = np.clip(y2, 0, orig_h)

        # NMS (非极大值抑制)
        boxes_for_nms = np.stack([x1, y1, x2, y2], axis=1)
        indices = cv2.dnn.NMSBoxes(
            boxes_for_nms.tolist(),
            confidences.tolist(),
            threshold,
            0.45,  # NMS IoU 阈值
        )

        detections = []

        # 处理 NMSBoxes 返回值的不同格式
        # OpenCV 不同版本返回格式不同：可能是 list、tuple、ndarray
        if indices is None or len(indices) == 0:
            return detections

        # 将 indices 转换为一维列表
        if isinstance(indices, np.ndarray):
            indices = indices.flatten().tolist()
        elif isinstance(indices, tuple):
            indices = list(indices)

        for idx in indices:
            # 确保 idx 是整数
            idx = int(idx) if not isinstance(idx, int) else idx

            cls_id = int(class_ids[idx])
            label = self.names.get(cls_id, f"class_{cls_id}")

            # 标签过滤
            if label_filter and label not in label_filter:
                continue

            x_min_val = int(x1[idx])
            y_min_val = int(y1[idx])
            x_max_val = int(x2[idx])
            y_max_val = int(y2[idx])
            area = (x_max_val - x_min_val) * (y_max_val - y_min_val)

            detections.append(
                {
                    "label": label,
                    "confidence": float(confidences[idx]),
                    "box": {
                        "x_min": x_min_val,
                        "y_min": y_min_val,
                        "x_max": x_max_val,
                        "y_max": y_max_val,
                    },
                    "area": area,
                    "norm_box": {
                        "x": (x_min_val + x_max_val) / 2 / orig_w,
                        "y": (y_min_val + y_max_val) / 2 / orig_h,
                        "w": (x_max_val - x_min_val) / orig_w,
                        "h": (y_max_val - y_min_val) / orig_h,
                    },
                }
            )

        return detections

    def detect(
        self,
        image: np.ndarray,
        threshold: float = 0.5,
        label_filter: list[str] | None = None,
        regions: list[tuple[int, int, int, int]] | None = None,
    ) -> tuple[list[dict], float]:
        """执行目标检测"""
        if not self.is_ready():
            raise RuntimeError("模型未加载")

        start_time = time.time()
        detections = []

        if regions and len(regions) > 0:
            for region in regions:
                x_min, y_min, x_max, y_max = region
                h, w = image.shape[:2]
                x_min = max(0, x_min)
                y_min = max(0, y_min)
                x_max = min(w, x_max)
                y_max = min(h, y_max)

                if x_max <= x_min or y_max <= y_min:
                    continue

                cropped = image[y_min:y_max, x_min:x_max]
                if cropped.size == 0:
                    continue

                region_detections = self._detect_single(
                    cropped, threshold, label_filter
                )

                for det in region_detections:
                    det["box"]["x_min"] += x_min
                    det["box"]["y_min"] += y_min
                    det["box"]["x_max"] += x_min
                    det["box"]["y_max"] += y_min
                    detections.append(det)
        else:
            detections = self._detect_single(image, threshold, label_filter)

        inference_time_ms = (time.time() - start_time) * 1000
        return detections, inference_time_ms

    def _detect_single(
        self, image: np.ndarray, threshold: float, label_filter: list[str] | None = None
    ) -> list[dict[str, Any]]:
        """对单张图像执行检测"""
        if not self.backend or not self.backend.is_ready():
            return []

        # 预处理
        input_tensor = self._preprocess(image)

        # 推理
        output = self.backend.infer(input_tensor)

        # 后处理
        original_shape = image.shape[:2]
        return self._postprocess(output, original_shape, threshold, label_filter)


class MotionDetector:
    """
    运动检测器 - 基于背景差分法
    用于在目标检测前预筛选有运动的帧，减少不必要的 AI 推理
    """

    def __init__(self):
        self.backgrounds: dict[str, np.ndarray] = {}
        self.motion_threshold = 25
        self.min_contour_area = 500

    def detect(
        self,
        image: np.ndarray,
        camera_name: str,
        zones: list[dict] | None = None,
    ) -> tuple[list[dict[str, Any]], bool]:
        h, w = image.shape[:2]
        if len(image.shape) == 3:
            gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        else:
            gray = image.copy()

        # 高斯模糊平滑噪点
        gray = cv2.GaussianBlur(gray, (21, 21), 0)

        if camera_name not in self.backgrounds:
            self.backgrounds[camera_name] = gray.astype(np.float32)
            return [], False

        cv2.accumulateWeighted(gray, self.backgrounds[camera_name], 0.1)

        frame_delta = cv2.absdiff(
            gray, cv2.convertScaleAbs(self.backgrounds[camera_name])
        )
        thresh = cv2.threshold(
            frame_delta, self.motion_threshold, 255, cv2.THRESH_BINARY
        )[1]

        # 多区域联合掩码：有 zones 时仅在所有区域内检测运动
        if zones:
            mask = np.zeros((h, w), dtype=np.uint8)
            for zone in zones:
                pts_flat = zone.get("points", [])
                if len(pts_flat) >= 6:  # 至少 3 个点
                    pts = [(int(pts_flat[i] * w), int(pts_flat[i + 1] * h)) for i in range(0, len(pts_flat), 2)]
                    pts_np = np.array([pts], dtype=np.int32)
                    cv2.fillPoly(mask, [pts_np], (255,))  # type: ignore
            thresh = cv2.bitwise_and(thresh, thresh, mask=mask)

        kernel = np.ones((3, 3), np.uint8)
        thresh = cv2.dilate(thresh, kernel, iterations=2)

        contours, _ = cv2.findContours(
            thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE
        )

        motion_boxes = []
        for contour in contours:
            if cv2.contourArea(contour) < self.min_contour_area:
                continue
            x, y, cw, ch = cv2.boundingRect(contour)
            motion_boxes.append(
                {"y_min": y, "x_min": x, "y_max": y + ch, "x_max": x + cw}
            )

        has_motion = len(motion_boxes) > 0
        return motion_boxes, has_motion
