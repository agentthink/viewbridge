package rtspadapter

import (
	"context"
	"log/slog"

	"github.com/gowvp/owl/internal/core/ipc"
	"github.com/gowvp/owl/internal/core/sms"
)

var _ ipc.Protocoler = (*Adapter)(nil)

// Adapter RTSP 协议适配器
// 处理 RTSP 拉流的状态管理
type Adapter struct {
	ipcCore ipc.Core
	smsCore sms.Core
}

// DeleteDevice implements ipc.Protocoler.
func (a *Adapter) DeleteDevice(ctx context.Context, device *ipc.Device) error {
	return nil
}

func NewAdapter(ipcCore ipc.Core, smsCore sms.Core) *Adapter {
	return &Adapter{
		ipcCore: ipcCore,
		smsCore: smsCore,
	}
}

// InitDevice implements ipc.Protocoler.
func (a *Adapter) InitDevice(ctx context.Context, device *ipc.Device) error {
	return nil
}

// OnStreamChanged implements ipc.Protocoler.
// RTSP 拉流断开时更新通道状态（IsOnline=false, IsPlaying=false）
func (a *Adapter) OnStreamChanged(ctx context.Context, app, stream string) error {
	// 通过 app+stream 查询通道，支持自定义 app/stream
	ch, err := a.ipcCore.GetChannelByAppStreamOrID(ctx, app, stream)
	if err != nil {
		slog.WarnContext(ctx, "RTSP 通道未找到", "app", app, "stream", stream, "err", err)
		return nil
	}
	if _, err := a.ipcCore.UpdateChannelOnlineAndPlaying(ctx, ch.Stream, false, false); err != nil {
		slog.WarnContext(ctx, "更新 RTSP 通道状态失败", "app", app, "stream", stream, "err", err)
	}
	return nil
}

// OnStreamNotFound implements ipc.Protocoler.
// 当流不存在时，从 Channel 获取配置并启动拉流代理
func (a *Adapter) OnStreamNotFound(ctx context.Context, app string, stream string) error {
	// 通过 app+stream 查询通道，支持自定义 app/stream
	ch, err := a.ipcCore.GetChannelByAppStreamOrID(ctx, app, stream)
	if err != nil {
		return err
	}

	svr, err := a.smsCore.GetMediaServer(ctx, sms.DefaultMediaServerID)
	if err != nil {
		return err
	}
	resp, err := a.smsCore.CreateStreamProxy(svr, sms.AddStreamProxyRequest{
		App:     ch.App,
		Stream:  ch.Stream,
		URL:     ch.Config.SourceURL,
		RTPType: ch.Config.Transport,
	})
	if err != nil {
		return err
	}

	// 更新 StreamKey 和 IsOnline（用于后续关闭拉流代理）
	_, err = a.ipcCore.UpdateChannelConfigAndOnline(ctx, ch.ID, true, func(cfg *ipc.StreamConfig) {
		cfg.StreamKey = resp.Data.Key
	})

	return err
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
// RTSP 协议不支持云台控制
func (a *Adapter) PTZControl(ctx context.Context, device *ipc.Device, channel *ipc.Channel, cmd ipc.PTZCommand) error {
	return nil
}
