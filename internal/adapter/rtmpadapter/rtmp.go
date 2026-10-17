package rtmpadapter

import (
	"context"
	"log/slog"

	"github.com/gowvp/owl/internal/conf"
	"github.com/gowvp/owl/internal/core/ipc"
	"github.com/ixugo/goddd/pkg/hook"
	"github.com/ixugo/goddd/pkg/orm"
)

var _ ipc.Protocoler = (*Adapter)(nil)

// Adapter RTMP 协议适配器
// 处理 RTMP 推流的鉴权、状态管理等逻辑
type Adapter struct {
	ipcCore ipc.Core
	conf    *conf.Bootstrap
}

func NewAdapter(ipcCore ipc.Core, conf *conf.Bootstrap) *Adapter {
	return &Adapter{
		ipcCore: ipcCore,
		conf:    conf,
	}
}

// DeleteDevice implements ipc.Protocoler.
func (a *Adapter) DeleteDevice(ctx context.Context, device *ipc.Device) error {
	return nil
}

// InitDevice implements ipc.Protocoler.
func (a *Adapter) InitDevice(ctx context.Context, device *ipc.Device) error {
	return nil
}

// OnStreamChanged implements ipc.Protocoler.
// RTMP 推流断开时更新通道状态（IsOnline=false, IsPlaying=false）
func (a *Adapter) OnStreamChanged(ctx context.Context, app, stream string) error {
	now := orm.Now()
	// 通过 app+stream 查询通道，支持自定义 app/stream
	ch, err := a.ipcCore.GetChannelByAppStreamOrID(ctx, app, stream)
	if err != nil {
		slog.WarnContext(ctx, "RTMP 通道未找到", "app", app, "stream", stream, "err", err)
		return nil
	}
	_, err = a.ipcCore.UpdateChannelConfigAndOnline(ctx, ch.ID, false, func(cfg *ipc.StreamConfig) {
		cfg.StoppedAt = &now
	})
	if err != nil {
		slog.WarnContext(ctx, "更新 RTMP 通道停流状态失败", "app", app, "stream", stream, "err", err)
	}
	// 同时更新 IsPlaying
	if _, err := a.ipcCore.UpdateChannelPlaying(ctx, ch.Stream, false); err != nil {
		slog.WarnContext(ctx, "更新 RTMP 通道播放状态失败", "stream", stream, "err", err)
	}
	return nil
}

// OnStreamNotFound implements ipc.Protocoler.
// RTMP 推流不需要处理流不存在事件（推流方需要主动推流）
func (a *Adapter) OnStreamNotFound(ctx context.Context, app string, stream string) error {
	return nil
}

// OnPublish 处理 RTMP 推流鉴权
// 验证推流参数中的 sign 字段是否与配置的 RTMPSecret MD5 一致
func (a *Adapter) OnPublish(ctx context.Context, app, stream string, params map[string]string) (bool, error) {
	// 通过 app+stream 查询通道，支持自定义 app/stream
	ch, err := a.ipcCore.GetChannelByAppStreamOrID(ctx, app, stream)
	if err != nil {
		return false, err
	}

	// 如果开启了鉴权
	if !ch.Config.IsAuthDisabled {
		sign := params["sign"]
		expectedSign := hook.MD5(a.conf.Server.RTMPSecret)
		if sign != expectedSign {
			return false, nil
		}
	}

	// 更新通道推流状态
	now := orm.Now()
	_, err = a.ipcCore.UpdateChannelConfigAndOnline(ctx, ch.ID, true, func(cfg *ipc.StreamConfig) {
		cfg.PushedAt = &now
		cfg.Session = params["session"]
	})
	if err != nil {
		slog.ErrorContext(ctx, "更新 RTMP 通道推流状态失败", "err", err)
	}

	return true, nil
}

// QueryCatalog implements ipc.Protocoler.
func (a *Adapter) QueryCatalog(ctx context.Context, device *ipc.Device) error {
	return nil
}

// StartPlay implements ipc.Protocoler.
func (a *Adapter) StartPlay(ctx context.Context, device *ipc.Device, channel *ipc.Channel) (*ipc.PlayResponse, error) {
	return nil, nil
}

// StopPlay implements ipc.Protocoler.
func (a *Adapter) StopPlay(ctx context.Context, device *ipc.Device, channel *ipc.Channel) error {
	return nil
}

// ValidateDevice implements ipc.Protocoler.
func (a *Adapter) ValidateDevice(ctx context.Context, device *ipc.Device) error {
	return nil
}

// PTZControl implements ipc.Protocoler.
// RTMP 协议不支持云台控制
func (a *Adapter) PTZControl(ctx context.Context, device *ipc.Device, channel *ipc.Channel, cmd ipc.PTZCommand) error {
	return nil
}
