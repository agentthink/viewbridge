package ipc

import (
	"context"
)

// Protocoler 协议抽象接口（端口）
//
// 设计原则:
// 1. 接口在 ipc 包内定义，避免循环依赖
// 2. 接口方法直接使用领域模型 (*Device, *Channel)
// 3. 适配器实现此接口，可以直接依赖和修改领域模型
// 4. 符合依赖倒置原则 (DIP):
//   - ipc (高层) 依赖 Protocoler 接口
//   - adapter (低层) 实现 Protocoler 接口
//   - adapter (低层) 依赖 ipc.Device (高层) ✅ 合理
//
// 这就是依赖反转！
type Protocoler interface {
	// ValidateDevice 验证设备连接（添加设备前调用）
	// 可以修改设备信息（如从 ONVIF 获取的固件版本等）
	ValidateDevice(ctx context.Context, device *Device) error

	// InitDevice 初始化设备连接（添加设备后调用）
	// 例如: GB28181 不需要主动初始化，ONVIF 需要查询 Profiles 作为通道
	InitDevice(ctx context.Context, device *Device) error

	// QueryCatalog 查询设备目录/通道
	QueryCatalog(ctx context.Context, device *Device) error

	// StartPlay 开始播放
	StartPlay(ctx context.Context, device *Device, channel *Channel) (*PlayResponse, error)

	// StopPlay 停止播放
	StopPlay(ctx context.Context, device *Device, channel *Channel) error

	DeleteDevice(ctx context.Context, device *Device) error

	// PTZControl 云台控制
	PTZControl(ctx context.Context, device *Device, channel *Channel, cmd PTZCommand) error

	Hooker
}

type Hooker interface {
	OnStreamNotFound(ctx context.Context, app, stream string) error
	// OnStreamChanged 流注销时调用，用于更新通道状态
	// app/stream 用于支持自定义 app/stream 的 RTMP/RTSP 通道
	OnStreamChanged(ctx context.Context, app, stream string) error
}

// OnPublisher 推流鉴权接口（可选实现）
// 只有 RTMP 需要实现此接口
type OnPublisher interface {
	// OnPublish 处理推流鉴权
	// 返回 true 表示鉴权通过，false 表示鉴权失败
	// app/stream 用于支持自定义 app/stream 的 RTMP/RTSP 通道
	OnPublish(ctx context.Context, app, stream string, params map[string]string) (bool, error)
}

// PlayResponse 播放响应
type PlayResponse struct {
	SSRC   string // GB28181 SSRC
	Stream string // 流 ID
	RTSP   string // RTSP 地址 (ONVIF)
}

// PTZCommand 云台控制命令
type PTZCommand struct {
	Action    string  `json:"action"`    // 动作类型: continuous(连续移动), absolute(绝对移动), relative(相对移动), stop(停止), preset(预置位)
	Direction string  `json:"direction"` // 方向: up, down, left, right, upleft, upright, downleft, downright, zoomin, zoomout
	Speed     float64 `json:"speed"`     // 速度 (0-1), 默认 0.5
	X         float64 `json:"x"`         // X 轴位置 (绝对/相对移动时使用, -1 到 1)
	Y         float64 `json:"y"`         // Y 轴位置 (绝对/相对移动时使用, -1 到 1)
	Zoom      float64 `json:"zoom"`      // 缩放值 (绝对/相对移动时使用, 0 到 1)
	PresetID  string  `json:"preset_id"` // 预置位 ID
	PresetOp  string  `json:"preset_op"` // 预置位操作: goto(转到), set(设置), remove(删除)
}
