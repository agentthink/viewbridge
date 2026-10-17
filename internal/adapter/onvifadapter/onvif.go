package onvifadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gowvp/onvif"
	devicemodel "github.com/gowvp/onvif/device"
	m "github.com/gowvp/onvif/media"
	ptzmodel "github.com/gowvp/onvif/ptz"
	sdkdevice "github.com/gowvp/onvif/sdk/device"
	sdkmedia "github.com/gowvp/onvif/sdk/media"
	sdkptz "github.com/gowvp/onvif/sdk/ptz"
	"github.com/gowvp/onvif/xsd"
	xsdonvif "github.com/gowvp/onvif/xsd/onvif"
	"github.com/gowvp/owl/internal/core/ipc"
	"github.com/gowvp/owl/internal/core/sms"
	"github.com/ixugo/goddd/pkg/conc"
	"github.com/ixugo/goddd/pkg/orm"
)

var _ ipc.Protocoler = (*Adapter)(nil)

// Adapter ONVIF 协议适配器
//
// 设计说明:
// - 适配器实现 ipc.Protocol 接口（Port 在 ipc 包内）
// - 适配器直接依赖领域模型 (ipc.Device, ipc.Channel)
// - 适配器依赖 ipc.Adapter 来访问存储和通用功能
// - 这符合清晰架构: 外层（适配器）依赖内层（领域）
type Adapter struct {
	devices conc.Map[string, *Device] // ONVIF 设备连接缓存
	adapter ipc.Adapter               // 通用适配器，提供 SaveChannels 等方法
	client  *http.Client
	sms     sms.Core
}

// Device ONVIF 设备包装（内存状态 + ONVIF 连接）
type Device struct {
	*onvif.Device
	KeepaliveAt orm.Time // 最后心跳时间
	IsOnline    bool     // 在线状态（内存缓存）
}

// DeleteDevice implements ipc.Protocoler.
func (a *Adapter) DeleteDevice(ctx context.Context, device *ipc.Device) error {
	a.devices.Delete(device.ID)
	return nil
}

func NewAdapter(adapter ipc.Adapter, sms sms.Core) *Adapter {
	cli := *http.DefaultClient
	cli.Timeout = time.Millisecond * 3000
	a := Adapter{
		adapter: adapter,
		client:  &cli,
		sms:     sms,
	}
	a.init()

	// 启动健康检查
	go a.startHealthCheck(context.Background())

	return &a
}

func (a *Adapter) init() {
	devices, err := a.adapter.ListDevices(context.TODO())
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		if device.IsOnvif() {
			go func(device *ipc.Device) {
				onvifDev, err := onvif.NewDevice(onvif.DeviceParams{
					Xaddr:      fmt.Sprintf("%s:%d", device.IP, device.Port),
					Username:   device.GetUsername(),
					Password:   device.Password,
					HttpClient: a.client,
				})
				if err != nil {
					_ = a.adapter.Update(device.ID, func(d *ipc.Device) {
						d.IsOnline = false
					})
					slog.Error("初始化 ONVIF 设备失败", "err", err, "device_id", device.ID)
				}
				if onvifDev == nil {
					return
				}
				a.devices.Store(device.ID, &Device{
					Device:   onvifDev,
					IsOnline: err == nil,
				})
			}(device)
		}
	}
}

// ValidateDevice 实现 ipc.Protocol 接口 - ONVIF 设备验证
func (a *Adapter) ValidateDevice(ctx context.Context, dev *ipc.Device) error {
	onvifDev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:      fmt.Sprintf("%s:%d", dev.IP, dev.Port),
		Username:   dev.GetUsername(),
		Password:   dev.Password,
		HttpClient: a.client,
	})
	if err != nil {
		return fmt.Errorf("IP 或 PORT 错误: %w", err)
	}

	// 获取设备信息并填充到领域模型
	resp, err := sdkdevice.Call_GetDeviceInformation(ctx, onvifDev, devicemodel.GetDeviceInformation{})
	if err != nil {
		return fmt.Errorf("账号或密码错误: %w", err)
	}
	dev.Transport = "tcp"
	dev.Ext.Firmware = resp.FirmwareVersion
	dev.Ext.Manufacturer = resp.Manufacturer
	dev.Ext.Model = resp.Model
	dev.IsOnline = true
	dev.Address = fmt.Sprintf("%s:%d", dev.IP, dev.Port)
	return nil
}

// InitDevice 实现 ipc.Protocol 接口 - 初始化 ONVIF 设备
// ONVIF 设备初始化时，自动查询 Profiles 并创建为通道
func (a *Adapter) InitDevice(ctx context.Context, dev *ipc.Device) error {
	// 创建 ONVIF 连接
	onvifDev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:      fmt.Sprintf("%s:%d", dev.IP, dev.Port),
		Username:   dev.GetUsername(),
		Password:   dev.Password,
		HttpClient: a.client,
	})
	if err != nil {
		return err
	}

	// 缓存设备连接
	d := Device{
		Device:   onvifDev,
		IsOnline: true,
	}
	a.devices.Store(dev.ID, &d)

	// 自动查询 Profiles 作为通道
	return a.queryAndSaveProfiles(ctx, dev, &d)
}

// QueryCatalog 实现 ipc.Protocol 接口 - ONVIF 查询 Profiles
func (a *Adapter) QueryCatalog(ctx context.Context, dev *ipc.Device) error {
	onvifDev, ok := a.devices.Load(dev.ID)
	if !ok {
		// 设备连接不在缓存中，尝试重新连接
		var err error
		d, err := onvif.NewDevice(onvif.DeviceParams{
			Xaddr:    fmt.Sprintf("%s:%d", dev.IP, dev.Port),
			Username: dev.GetUsername(),
			Password: dev.Password,
		})
		if err != nil {
			return fmt.Errorf("ONVIF 设备未初始化: %w", err)
		}
		onvifDev = &Device{
			Device:   d,
			IsOnline: true,
		}
		a.devices.Store(dev.ID, onvifDev)
	}

	return a.queryAndSaveProfiles(ctx, dev, onvifDev)
}

// StartPlay 实现 ipc.Protocol 接口 - ONVIF 播放
func (a *Adapter) StartPlay(ctx context.Context, dev *ipc.Device, ch *ipc.Channel) (*ipc.PlayResponse, error) {
	onvifDev, ok := a.devices.Load(dev.ID)
	if !ok {
		return nil, fmt.Errorf("ONVIF 设备未初始化")
	}

	// 获取 RTSP 地址
	streamURI, err := a.getStreamURI(ctx, onvifDev, ch.ChannelID)
	if err != nil {
		return nil, err
	}

	return &ipc.PlayResponse{
		RTSP: streamURI,
	}, nil
}

// StopPlay 实现 ipc.Protocol 接口 - ONVIF 停止播放
func (a *Adapter) StopPlay(ctx context.Context, dev *ipc.Device, ch *ipc.Channel) error {
	// ONVIF 通常不需要显式停止播放
	return nil
}

// queryAndSaveProfiles 查询 ONVIF Profiles 并保存为通道
//
// 使用统一的 SaveChannels 方法，自动处理增量更新和删除
func (a *Adapter) queryAndSaveProfiles(ctx context.Context, device *ipc.Device, onvifDev *Device) error {
	resp, err := sdkmedia.Call_GetProfiles(ctx, onvifDev.Device, m.GetProfiles{})
	if err != nil {
		return fmt.Errorf("账号或密码错误: %w", err)
	}

	// 将 Profiles 转换为通道列表
	channels := make([]*ipc.Channel, 0, len(resp.Profiles))
	for _, profile := range resp.Profiles {
		channel := &ipc.Channel{
			DeviceID:  device.ID,
			ChannelID: string(profile.Token),
			Name:      string(profile.Name),
			DID:       device.ID,
			IsOnline:  true,
			Type:      ipc.TypeOnvif,
		}
		channels = append(channels, channel)
	}
	if len(channels) == 0 {
		return fmt.Errorf("没有找到 ONVIF 通道")
	}

	// 使用统一的 SaveChannels 方法保存（自动处理增删改）
	if err := a.adapter.SaveChannels(channels); err != nil {
		return fmt.Errorf("保存 ONVIF 通道失败: %w", err)
	}

	slog.InfoContext(ctx, "ONVIF Profiles 同步完成",
		"device_id", device.ID,
		"profile_count", len(channels))

	return nil
}

// getStreamURI 获取 RTSP 流地址
func (a *Adapter) getStreamURI(ctx context.Context, dev *Device, profileToken string) (string, error) {
	var param m.GetStreamUri
	param.StreamSetup.Transport.Protocol = "RTSP"
	param.StreamSetup.Stream = "RTP-Unicast"
	param.ProfileToken = xsdonvif.ReferenceToken(profileToken)
	resp, err := sdkmedia.Call_GetStreamUri(ctx, dev.Device, param)
	if err != nil {
		return "", err
	}
	playURL := buildPlayURL(string(resp.MediaUri.Uri), dev.Device.GetDeviceParams().Username, dev.Device.GetDeviceParams().Password)
	return playURL, nil
}

func buildPlayURL(rawurl, username, password string) string {
	if username != "" && password != "" {
		return strings.Replace(rawurl, "rtsp://", fmt.Sprintf("rtsp://%s:%s@", username, password), 1)
	}
	return rawurl
}

func (a *Adapter) Discover(ctx context.Context, w io.Writer) error {
	recv, err := onvif.AllAvailableDevicesAtSpecificEthernetInterfaces()
	if err != nil {
		return err
	}

	for {
		select {
		case dev, ok := <-recv:
			if !ok {
				return nil
			}
			var exists bool
			a.devices.Range(func(key string, value *Device) bool {
				if value.GetDeviceParams().Xaddr == dev.GetDeviceParams().Xaddr {
					exists = true
					return false
				}
				return true
			})
			if exists {
				continue
			}
			b, _ := json.Marshal(toDiscoverResponse(dev))
			_, _ = w.Write(b)
		case <-ctx.Done():
			return nil
		case <-time.After(3 * time.Second):
			slog.DebugContext(ctx, "discover timeout")
			return nil
		}
	}
}

// PTZControl 实现 ipc.Protocoler 接口 - ONVIF 云台控制
func (a *Adapter) PTZControl(ctx context.Context, dev *ipc.Device, ch *ipc.Channel, cmd ipc.PTZCommand) error {
	onvifDev, ok := a.devices.Load(dev.ID)
	if !ok {
		return fmt.Errorf("ONVIF 设备未初始化")
	}

	// 根据动作类型执行不同的 PTZ 操作
	switch cmd.Action {
	case "continuous":
		return a.continuousMove(ctx, onvifDev, ch, cmd)
	case "stop":
		return a.stopMove(ctx, onvifDev, ch, cmd)
	case "absolute":
		return a.absoluteMove(ctx, onvifDev, ch, cmd)
	case "relative":
		return a.relativeMove(ctx, onvifDev, ch, cmd)
	case "preset":
		return a.presetOperation(ctx, onvifDev, ch, cmd)
	default:
		return fmt.Errorf("不支持的 PTZ 动作类型: %s", cmd.Action)
	}
}

// continuousMove 连续移动
func (a *Adapter) continuousMove(ctx context.Context, onvifDev *Device, ch *ipc.Channel, cmd ipc.PTZCommand) error {
	speed := cmd.Speed
	if speed <= 0 || speed > 1 {
		speed = 0.5
	}

	// 构建速度向量
	velocity := xsdonvif.PTZSpeed{
		PanTilt: xsdonvif.Vector2D{},
		Zoom:    xsdonvif.Vector1D{},
	}

	// 根据方向设置 Pan/Tilt 速度
	switch cmd.Direction {
	case "up":
		velocity.PanTilt.Y = speed
	case "down":
		velocity.PanTilt.Y = -speed
	case "left":
		velocity.PanTilt.X = -speed
	case "right":
		velocity.PanTilt.X = speed
	case "upleft":
		velocity.PanTilt.X = -speed
		velocity.PanTilt.Y = speed
	case "upright":
		velocity.PanTilt.X = speed
		velocity.PanTilt.Y = speed
	case "downleft":
		velocity.PanTilt.X = -speed
		velocity.PanTilt.Y = -speed
	case "downright":
		velocity.PanTilt.X = speed
		velocity.PanTilt.Y = -speed
	case "zoomin":
		velocity.Zoom.X = speed
	case "zoomout":
		velocity.Zoom.X = -speed
	default:
		return fmt.Errorf("不支持的方向: %s", cmd.Direction)
	}

	req := ptzmodel.ContinuousMove{
		ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
		Velocity:     velocity,
	}

	_, err := sdkptz.Call_ContinuousMove(ctx, onvifDev.Device, req)
	if err != nil {
		return fmt.Errorf("ONVIF 连续移动失败: %w", err)
	}

	return nil
}

// stopMove 停止移动
func (a *Adapter) stopMove(ctx context.Context, onvifDev *Device, ch *ipc.Channel, cmd ipc.PTZCommand) error {
	req := ptzmodel.Stop{
		ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
		PanTilt:      true,
		Zoom:         true,
	}

	_, err := sdkptz.Call_Stop(ctx, onvifDev.Device, req)
	if err != nil {
		return fmt.Errorf("ONVIF 停止移动失败: %w", err)
	}

	return nil
}

// absoluteMove 绝对移动
func (a *Adapter) absoluteMove(ctx context.Context, onvifDev *Device, ch *ipc.Channel, cmd ipc.PTZCommand) error {
	position := xsdonvif.PTZVector{
		PanTilt: xsdonvif.Vector2D{},
		Zoom:    xsdonvif.Vector1D{},
	}

	// 设置位置
	if cmd.X != 0 || cmd.Y != 0 {
		position.PanTilt.X = cmd.X
		position.PanTilt.Y = cmd.Y
	}
	if cmd.Zoom != 0 {
		position.Zoom.X = cmd.Zoom
	}

	// 设置速度（可选）
	speed := xsdonvif.PTZSpeed{}
	if cmd.Speed > 0 && cmd.Speed <= 1 {
		speed.PanTilt.X = cmd.Speed
		speed.PanTilt.Y = cmd.Speed
		speed.Zoom.X = cmd.Speed
	}

	req := ptzmodel.AbsoluteMove{
		ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
		Position:     position,
		Speed:        speed,
	}

	_, err := sdkptz.Call_AbsoluteMove(ctx, onvifDev.Device, req)
	if err != nil {
		return fmt.Errorf("ONVIF 绝对移动失败: %w", err)
	}

	return nil
}

// relativeMove 相对移动
func (a *Adapter) relativeMove(ctx context.Context, onvifDev *Device, ch *ipc.Channel, cmd ipc.PTZCommand) error {
	translation := xsdonvif.PTZVector{
		PanTilt: xsdonvif.Vector2D{},
		Zoom:    xsdonvif.Vector1D{},
	}

	// 设置相对位移
	if cmd.X != 0 || cmd.Y != 0 {
		translation.PanTilt.X = cmd.X
		translation.PanTilt.Y = cmd.Y
	}
	if cmd.Zoom != 0 {
		translation.Zoom.X = cmd.Zoom
	}

	// 设置速度（可选）
	speed := xsdonvif.PTZSpeed{}
	if cmd.Speed > 0 && cmd.Speed <= 1 {
		speed.PanTilt.X = cmd.Speed
		speed.PanTilt.Y = cmd.Speed
		speed.Zoom.X = cmd.Speed
	}

	req := ptzmodel.RelativeMove{
		ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
		Translation:  translation,
		Speed:        speed,
	}

	_, err := sdkptz.Call_RelativeMove(ctx, onvifDev.Device, req)
	if err != nil {
		return fmt.Errorf("ONVIF 相对移动失败: %w", err)
	}

	return nil
}

// presetOperation 预置位操作
func (a *Adapter) presetOperation(ctx context.Context, onvifDev *Device, ch *ipc.Channel, cmd ipc.PTZCommand) error {
	if cmd.PresetID == "" {
		return fmt.Errorf("预置位 ID 不能为空")
	}

	switch cmd.PresetOp {
	case "goto":
		req := ptzmodel.GotoPreset{
			ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
			PresetToken:  xsdonvif.ReferenceToken(cmd.PresetID),
		}
		_, err := sdkptz.Call_GotoPreset(ctx, onvifDev.Device, req)
		if err != nil {
			return fmt.Errorf("ONVIF 转到预置位失败: %w", err)
		}
	case "set":
		req := ptzmodel.SetPreset{
			ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
			PresetName:   xsd.String(cmd.PresetID),
		}
		_, err := sdkptz.Call_SetPreset(ctx, onvifDev.Device, req)
		if err != nil {
			return fmt.Errorf("ONVIF 设置预置位失败: %w", err)
		}
	case "remove":
		req := ptzmodel.RemovePreset{
			ProfileToken: xsdonvif.ReferenceToken(ch.ChannelID),
			PresetToken:  xsdonvif.ReferenceToken(cmd.PresetID),
		}
		_, err := sdkptz.Call_RemovePreset(ctx, onvifDev.Device, req)
		if err != nil {
			return fmt.Errorf("ONVIF 删除预置位失败: %w", err)
		}
	default:
		return fmt.Errorf("不支持的预置位操作: %s", cmd.PresetOp)
	}

	return nil
}
