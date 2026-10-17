package api

import "github.com/ixugo/goddd/pkg/orm"

// AIKeepaliveInput 心跳回调请求体
type AIKeepaliveInput struct {
	Timestamp int64          `json:"timestamp"` // Unix 时间戳 (毫秒)
	Stats     *AIGlobalStats `json:"stats"`     // 全局统计信息
	Message   string         `json:"message"`   // 附加消息
}

// AIStartedInput 服务启动回调请求体
type AIStartedInput struct {
	Timestamp int64  `json:"timestamp"` // Unix 时间戳 (毫秒)
	Message   string `json:"message"`   // 启动消息
}

// AIDetectionInput 检测事件回调请求体
type AIDetectionInput struct {
	CameraID       string        `json:"camera_id"`       // 摄像头 ID
	Timestamp      orm.Time      `json:"timestamp"`       // Unix 时间戳 (毫秒)
	Detections     []AIDetection `json:"detections"`      // 检测结果列表
	Snapshot       string        `json:"snapshot"`        // Base64 编码的快照 (JPEG)
	SnapshotWidth  int           `json:"snapshot_width"`  // 快照宽度
	SnapshotHeight int           `json:"snapshot_height"` // 快照高度
}

// AIStoppedInput 任务停止回调请求体
type AIStoppedInput struct {
	CameraID  string   `json:"camera_id"` // 摄像头 ID
	Timestamp orm.Time `json:"timestamp"` // Unix 时间戳 (毫秒)
	Reason    string   `json:"reason"`    // 停止原因 (user_requested, error)
	Message   string   `json:"message"`   // 详细信息
}

// AIDetection 检测对象
type AIDetection struct {
	Label      string        `json:"label"`      // 物体类别
	Confidence float64       `json:"confidence"` // 置信度 (0.0 - 1.0)
	Box        AIBoundingBox `json:"box"`        // 像素坐标边界框
	Area       int           `json:"area"`       // 边界框像素面积
	NormBox    *AINormBox    `json:"norm_box"`   // 归一化边界框
}

// AIBoundingBox 像素坐标边界框
type AIBoundingBox struct {
	XMin int `json:"x_min"`
	YMin int `json:"y_min"`
	XMax int `json:"x_max"`
	YMax int `json:"y_max"`
}

// AINormBox 归一化边界框
type AINormBox struct {
	X float64 `json:"x"` // 中心点 X 坐标
	Y float64 `json:"y"` // 中心点 Y 坐标
	W float64 `json:"w"` // 宽度
	H float64 `json:"h"` // 高度
}

// AIGlobalStats 全局统计信息
type AIGlobalStats struct {
	ActiveStreams   int   `json:"active_streams"`   // 活跃流数量
	TotalDetections int64 `json:"total_detections"` // 总检测次数
	UptimeSeconds   int64 `json:"uptime_seconds"`   // 运行时间 (秒)
}

// AIWebhookOutput 通用响应体
type AIWebhookOutput struct {
	Code int    `json:"code"` // 错误代码，0 表示成功
	Msg  string `json:"msg"`  // 消息
}

func newAIWebhookOutputOK() AIWebhookOutput {
	return AIWebhookOutput{Code: 0, Msg: "success"}
}
