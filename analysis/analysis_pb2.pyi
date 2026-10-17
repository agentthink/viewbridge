from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class StartCameraRequest(_message.Message):
    __slots__ = ("camera_id", "camera_name", "rtsp_url", "detect_interval_seconds", "labels", "threshold", "zones", "alert_cooldown_seconds", "retry_limit", "callback_url", "callback_secret")
    CAMERA_ID_FIELD_NUMBER: _ClassVar[int]
    CAMERA_NAME_FIELD_NUMBER: _ClassVar[int]
    RTSP_URL_FIELD_NUMBER: _ClassVar[int]
    DETECT_INTERVAL_SECONDS_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    ZONES_FIELD_NUMBER: _ClassVar[int]
    ALERT_COOLDOWN_SECONDS_FIELD_NUMBER: _ClassVar[int]
    RETRY_LIMIT_FIELD_NUMBER: _ClassVar[int]
    CALLBACK_URL_FIELD_NUMBER: _ClassVar[int]
    CALLBACK_SECRET_FIELD_NUMBER: _ClassVar[int]
    camera_id: str
    camera_name: str
    rtsp_url: str
    detect_interval_seconds: float
    labels: _containers.RepeatedScalarFieldContainer[str]
    threshold: float
    zones: _containers.RepeatedCompositeFieldContainer[AnalysisZone]
    alert_cooldown_seconds: float
    retry_limit: int
    callback_url: str
    callback_secret: str
    def __init__(self, camera_id: _Optional[str] = ..., camera_name: _Optional[str] = ..., rtsp_url: _Optional[str] = ..., detect_interval_seconds: _Optional[float] = ..., labels: _Optional[_Iterable[str]] = ..., threshold: _Optional[float] = ..., zones: _Optional[_Iterable[_Union[AnalysisZone, _Mapping]]] = ..., alert_cooldown_seconds: _Optional[float] = ..., retry_limit: _Optional[int] = ..., callback_url: _Optional[str] = ..., callback_secret: _Optional[str] = ...) -> None: ...

class StartCameraResponse(_message.Message):
    __slots__ = ("success", "message", "source_width", "source_height", "source_fps")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_WIDTH_FIELD_NUMBER: _ClassVar[int]
    SOURCE_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FPS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    message: str
    source_width: int
    source_height: int
    source_fps: float
    def __init__(self, success: bool = ..., message: _Optional[str] = ..., source_width: _Optional[int] = ..., source_height: _Optional[int] = ..., source_fps: _Optional[float] = ...) -> None: ...

class StopCameraRequest(_message.Message):
    __slots__ = ("camera_id",)
    CAMERA_ID_FIELD_NUMBER: _ClassVar[int]
    camera_id: str
    def __init__(self, camera_id: _Optional[str] = ...) -> None: ...

class StopCameraResponse(_message.Message):
    __slots__ = ("success", "message")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    success: bool
    message: str
    def __init__(self, success: bool = ..., message: _Optional[str] = ...) -> None: ...

class StatusRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class StatusResponse(_message.Message):
    __slots__ = ("is_ready", "cameras", "stats")
    IS_READY_FIELD_NUMBER: _ClassVar[int]
    CAMERAS_FIELD_NUMBER: _ClassVar[int]
    STATS_FIELD_NUMBER: _ClassVar[int]
    is_ready: bool
    cameras: _containers.RepeatedCompositeFieldContainer[CameraStatus]
    stats: GlobalStats
    def __init__(self, is_ready: bool = ..., cameras: _Optional[_Iterable[_Union[CameraStatus, _Mapping]]] = ..., stats: _Optional[_Union[GlobalStats, _Mapping]] = ...) -> None: ...

class CameraStatus(_message.Message):
    __slots__ = ("camera_id", "status", "frames_processed", "last_error", "retry_count")
    CAMERA_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    FRAMES_PROCESSED_FIELD_NUMBER: _ClassVar[int]
    LAST_ERROR_FIELD_NUMBER: _ClassVar[int]
    RETRY_COUNT_FIELD_NUMBER: _ClassVar[int]
    camera_id: str
    status: str
    frames_processed: int
    last_error: str
    retry_count: int
    def __init__(self, camera_id: _Optional[str] = ..., status: _Optional[str] = ..., frames_processed: _Optional[int] = ..., last_error: _Optional[str] = ..., retry_count: _Optional[int] = ...) -> None: ...

class GlobalStats(_message.Message):
    __slots__ = ("active_streams", "total_detections", "uptime_seconds")
    ACTIVE_STREAMS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_DETECTIONS_FIELD_NUMBER: _ClassVar[int]
    UPTIME_SECONDS_FIELD_NUMBER: _ClassVar[int]
    active_streams: int
    total_detections: int
    uptime_seconds: int
    def __init__(self, active_streams: _Optional[int] = ..., total_detections: _Optional[int] = ..., uptime_seconds: _Optional[int] = ...) -> None: ...

class HealthCheckRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class HealthCheckResponse(_message.Message):
    __slots__ = ("status",)
    class ServingStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[HealthCheckResponse.ServingStatus]
        SERVING: _ClassVar[HealthCheckResponse.ServingStatus]
        NOT_SERVING: _ClassVar[HealthCheckResponse.ServingStatus]
    UNKNOWN: HealthCheckResponse.ServingStatus
    SERVING: HealthCheckResponse.ServingStatus
    NOT_SERVING: HealthCheckResponse.ServingStatus
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: HealthCheckResponse.ServingStatus
    def __init__(self, status: _Optional[_Union[HealthCheckResponse.ServingStatus, str]] = ...) -> None: ...

class AnalysisZone(_message.Message):
    __slots__ = ("points", "labels", "name")
    POINTS_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    points: _containers.RepeatedScalarFieldContainer[float]
    labels: _containers.RepeatedScalarFieldContainer[str]
    name: str
    def __init__(self, points: _Optional[_Iterable[float]] = ..., labels: _Optional[_Iterable[str]] = ..., name: _Optional[str] = ...) -> None: ...
