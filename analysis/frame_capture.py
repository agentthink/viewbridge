from collections import deque
import logging
import os
import queue
import subprocess
import threading
import time
from typing import Deque, Optional
import numpy as np


slog = logging.getLogger("Capture")


class LogPipe(threading.Thread):
    def __init__(self, log_name: str):
        super().__init__(daemon=True)
        self.logger = logging.getLogger(log_name)
        self.deque: Deque[str] = deque(maxlen=100)
        self.fd_read, self.fd_write = os.pipe()
        self.pipe_reader = os.fdopen(self.fd_read)
        self._closed = False
        self.start()

    def fileno(self):
        return self.fd_write

    def run(self):
        # 使用 iter() 包装 self.pipe_reader.readline 方法和空字符串""作为哨兵，使其不断读取管道内容。
        # iter(self.pipe_reader.readline, "") 会不断调用 readline()，直到返回空字符串（代表 EOF），循环终止。
        try:
            for line in iter(self.pipe_reader.readline, ""):
                self.deque.append(line)
        except (OSError, ValueError):
            # 管道已关闭，忽略错误
            pass
        finally:
            try:
                if not self._closed:
                    self.pipe_reader.close()
            except (OSError, ValueError):
                pass

    def dump(self):
        while len(self.deque) > 0:
            self.logger.error(self.deque.popleft())

    def close(self):
        # 先关闭写端，让读端线程收到 EOF 并退出
        if self._closed:
            return
        self._closed = True
        try:
            os.close(self.fd_write)
        except OSError:
            pass
        # 等待读线程结束
        self.join(timeout=1)


class FrameCapture:
    def __init__(
        self,
        rtsp_url: str,
        output_queue: queue.Queue,
        detect_interval: float = 5.0,
        retry_limit: int = 10,
    ):
        self.rtsp_url = rtsp_url
        self.output_queue = output_queue
        # 将间隔秒数转为 ffmpeg 能接受的分数帧率字符串，如 2.0 -> "1/2"、0.5 -> "2/1"
        self.target_fps_str = self._interval_to_fps_str(detect_interval)
        self.retry_limit = retry_limit
        self._stop_event = threading.Event()
        self._thread: Optional[threading.Thread] = None
        self._proccess: Optional[subprocess.Popen] = None

        # 流信息
        self.width = 0
        self.height = 0
        self.fps = 0.0

        # 错误状态，供外部查询
        self.error_count = 0
        self.last_error = ""
        self.is_failed = False

    @staticmethod
    def _interval_to_fps_str(interval: float) -> str:
        """
        把分析间隔（秒）转成 ffmpeg -r 接受的分数帧率字符串。
        interval=2.0 → "1/2"，interval=0.5 → "2/1"，interval=0.1 → "10/1"。
        使用整数分数避免浮点精度问题。
        """
        from fractions import Fraction
        frac = Fraction(1 / interval).limit_denominator(1000)
        return f"{frac.numerator}/{frac.denominator}"

    def start(self):
        if self._thread is not None and self._thread.is_alive():
            return
        self._stop_event.clear()
        self._thread = threading.Thread(target=self._capture_loop, daemon=True)
        self._thread.start()
        slog.info(f"FrameCapture started for {self.rtsp_url}")

    def stop(self):
        # 设置停止事件
        self._stop_event.set()
        # 终止进程
        self._terminate_process()
        # 等待线程结束
        if self._thread:
            self._thread.join(timeout=2)
        slog.info(f"FrameCapture stopped for {self.rtsp_url}")

    def _get_stream_info(self) -> bool:
        slog.debug(f"正在探测流信息... {self.rtsp_url}")
        ffprobe_cmd = [
            "ffprobe",
            "-v",
            "error",
            "-select_streams",
            "v:0",
            "-show_entries",
            "stream=width,height,r_frame_rate",
            "-of",
            "csv=p=0",
            "-rtsp_transport",
            "tcp",  # 强制 TCP 更稳定
            self.rtsp_url,
        ]
        try:
            # 执行一个外部命令（比如系统命令、shell 脚本、其他可执行程序），并直接返回该命令在标准输出（stdout）中打印的内容。
            output = (
                subprocess.check_output(ffprobe_cmd, timeout=15).decode("utf-8").strip()
            )
            parts = output.split(",")
            if len(parts) >= 2:
                self.width = int(parts[0])
                self.height = int(parts[1])
                if len(parts) >= 3 and "/" in parts[2]:
                    num, den = parts[2].split("/")
                    self.fps = float(num) / float(den)
                else:
                    self.fps = 25.0
                slog.info(
                    f"ffprobe 探测成功: {self.width}x{self.height} @ {self.fps:.2f}fps"
                )
                return True
        except Exception as e:
            slog.error(f"探测流信息失败: {e}")

        return False

    def _capture_loop(self):

        log_pipe: Optional[LogPipe] = None
        while not self._stop_event.is_set():
            if self.width == 0 or self.height == 0:
                if not self._get_stream_info():
                    self.error_count += 1
                    if self.error_count >= self.retry_limit:
                        self.is_failed = True
                        self.last_error = (
                            f"探测流信息失败，已重试 {self.error_count} 次"
                        )
                        slog.error(self.last_error)
                        return
                    time.sleep(3)
                    continue
            # 成功获取流信息后重置错误计数
            self.error_count = 0
            if log_pipe:
                log_pipe.close()
            log_pipe = LogPipe(f"ffmpeg.{self.rtsp_url}")
            ffmpeg_cmd = [
                "ffmpeg",
                "-hide_banner",
                "-loglevel",
                "warning",  # 只输出 warning 以上，减少 IO
                "-rtsp_transport",
                "tcp",
                "-i",
                self.rtsp_url,
                "-f",
                "rawvideo",
                "-pix_fmt",
                "bgr24",  # 直接输出 OpenCV 友好的 BGR 格式
                "-r",
                self.target_fps_str,  # 分析帧率（分数形式，如 1/2 表示每 2 秒 1 帧）
                "pipe:1",
            ]
            slog.info(f"启动 ffmpeg 进程: {' '.join(ffmpeg_cmd[:-2])} ...")

            try:
                self._proccess = subprocess.Popen(
                    ffmpeg_cmd,
                    stdout=subprocess.PIPE,
                    stderr=log_pipe.fileno(),
                    bufsize=10**7,
                )
            except Exception as e:
                slog.error(f"启动 ffmpeg 进程失败: {e}")
                if log_pipe:
                    log_pipe.dump()
                time.sleep(3)
                continue
            frame_size = self.width * self.height * 3
            slog.info(f"开始读取帧 (size={frame_size})...")

            while not self._stop_event.is_set():
                try:
                    if self._proccess.poll() is not None:
                        slog.error("FFmpeg 进程意外退出")
                        log_pipe.dump()
                        break
                    if self._proccess.stdout is None:
                        slog.error("FFmpeg 进程 stdout 为空")
                        break
                    raw_bytes = self._proccess.stdout.read(frame_size)

                    if len(raw_bytes) != frame_size:
                        slog.warning("读取到不完整的帧 (流中断?)")
                        log_pipe.dump()  # 可能有网络错误
                        break
                    image = np.frombuffer(raw_bytes, dtype=np.uint8).reshape(
                        self.height, self.width, 3
                    )

                    try:
                        while not self.output_queue.empty():
                            try:
                                self.output_queue.get_nowait()
                            except queue.Empty:
                                break
                        self.output_queue.put_nowait(image)
                    except Exception as e:
                        pass

                except Exception as e:
                    slog.error(f"读取帧失败: {e}")
                    self.last_error = str(e)
                    log_pipe.dump()
                    break
            self._terminate_process()

            if log_pipe:
                log_pipe.close()

            if self._stop_event.is_set():
                break

            # ffmpeg 进程异常退出也计入错误计数
            self.error_count += 1
            if self.error_count >= self.retry_limit:
                self.is_failed = True
                self.last_error = f"帧捕获失败，已重试 {self.error_count} 次"
                slog.error(self.last_error)
                return
            time.sleep(2)

    def _terminate_process(self):
        if self._proccess:
            if self._proccess.poll() is None:
                self._proccess.terminate()
                try:
                    self._proccess.wait(timeout=2)
                except subprocess.TimeoutExpired:
                    self._proccess.kill()
            self._proccess = None

    def _hide_password(self, url):
        """隐藏 URL 中的密码"""
        try:
            if "@" in url:
                parts = url.split("@")
                if "//" in parts[0]:
                    protocol_auth = parts[0].split("//")
                    if ":" in protocol_auth[1]:
                        user = protocol_auth[1].split(":")[0]
                        return f"{protocol_auth[0]}//{user}:***@{parts[1]}"
            return url
        except:
            return url

    def get_stream_info(self):
        """返回流的基本信息"""
        return self.width, self.height, self.fps
