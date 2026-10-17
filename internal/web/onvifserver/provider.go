package onvifserver

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gowvp/onvif/server"
	"github.com/gowvp/owl/internal/conf"
	"github.com/gowvp/owl/internal/core/ipc"
	"github.com/gowvp/owl/internal/core/sms"
	"github.com/ixugo/goddd/pkg/web"
)

// StreamProvider 将平台通道映射为 ONVIF Profile，并生成 ZLM RTSP 地址。
type StreamProvider struct {
	ipc     ipc.Core
	smsCore sms.Core
	conf    *conf.Bootstrap

	mu      sync.RWMutex
	channels []*ipc.Channel
	devices  map[string]deviceMeta // DID -> 父设备展示名与在线状态
}

type deviceMeta struct {
	name   string
	online bool
}

var _ server.StreamProvider = (*StreamProvider)(nil)

// NewStreamProvider 创建 ONVIF 流提供者，启动时预加载通道列表。
func NewStreamProvider(ipcCore ipc.Core, smsCore sms.Core, cfg *conf.Bootstrap) *StreamProvider {
	p := &StreamProvider{ipc: ipcCore, smsCore: smsCore, conf: cfg}
	p.refreshChannels()
	return p
}

// refreshChannels 从数据库拉取全部通道与设备名，供 ONVIF Profile 列表使用。
func (p *StreamProvider) refreshChannels() {
	ctx := context.Background()
	channels, _, err := p.ipc.ListChannels(ctx, &ipc.FindChannelInput{
		PagerFilter: web.NewPagerFilterMaxSize(),
	})
	if err != nil {
		slog.Warn("onvif refresh channels failed", "err", err)
		return
	}
	byDID := make(map[string]deviceMeta)
	devs, _, err := p.ipc.ListDevices(ctx, &ipc.FindDeviceInput{
		PagerFilter: web.NewPagerFilterMaxSize(),
	})
	if err != nil {
		slog.Warn("onvif refresh devices failed", "err", err)
	} else {
		for _, dev := range devs {
			byDID[dev.ID] = deviceMeta{name: dev.GetName(), online: dev.IsOnline}
		}
	}
	p.mu.Lock()
	p.channels = channels
	p.devices = byDID
	p.mu.Unlock()
}

// ListStreamProfiles 返回全部通道 Profile；Token 用通道主键避免展示名重复导致 HA 少实体。
func (p *StreamProvider) ListStreamProfiles() []server.StreamProfile {
	p.refreshChannels()
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]server.StreamProfile, 0, len(p.channels))
	for _, ch := range p.channels {
		if ch.ID == "" {
			continue
		}
		out = append(out, server.StreamProfile{Token: ch.ID, Name: p.profileName(ch)})
	}
	slog.Debug("onvif list stream profiles", "count", len(out))
	return out
}

// GetStreamURI 根据 Profile Token 拼接可访问的 RTSP 地址。
func (p *StreamProvider) GetStreamURI(profileToken, requestHost string) string {
	ch := p.findChannel(profileToken)
	if ch == nil {
		slog.Warn("onvif stream uri: channel not found", "token", profileToken)
		return ""
	}
	if !ch.IsRTSP() && !ipc.ChannelReachable(ch, p.devices[ch.DID].online) {
		slog.Debug("onvif stream uri: channel unreachable", "token", profileToken, "channel_id", ch.ID)
		return ""
	}
	return p.buildRTSPURL(ch)
}

// GetSnapshotURI 暂不实现平台级截图，返回空由服务端填占位。
func (p *StreamProvider) GetSnapshotURI(profileToken, requestHost string) string {
	_ = profileToken
	_ = requestHost
	return ""
}

func (p *StreamProvider) findChannel(token string) *ipc.Channel {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, ch := range p.channels {
		if ch.ID == token {
			return ch
		}
	}
	return nil
}

func (p *StreamProvider) profileName(ch *ipc.Channel) string {
	dev := p.devices[ch.DID]
	devName := dev.name
	if devName == "" {
		devName = ch.DID
	}
	return ipc.ONVIFProfileName(ch, devName, dev.online)
}

// ProfileSnapshot 返回当前 ONVIF 暴露列表，仅供 GET /onvif/profiles 排障。
func (p *StreamProvider) ProfileSnapshot() []ProfileSnapshot {
	p.refreshChannels()
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]ProfileSnapshot, 0, len(p.channels))
	for _, ch := range p.channels {
		if ch.ID == "" {
			continue
		}
		dev := p.devices[ch.DID]
		out = append(out, ProfileSnapshot{
			Token:     ch.ID,
			Name:      p.profileName(ch),
			ChannelID: ch.ChannelID,
			DID:       ch.DID,
			Online:    ipc.ChannelReachable(ch, dev.online),
		})
	}
	return out
}

// ProfileSnapshot 单条 ONVIF Profile 调试信息。
type ProfileSnapshot struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	ChannelID string `json:"channel_id"`
	DID       string `json:"did"`
	Online    bool   `json:"online"`
}

func (p *StreamProvider) buildRTSPURL(ch *ipc.Channel) string {
	mediaID := ch.Config.MediaServerID
	if mediaID == "" {
		mediaID = sms.DefaultMediaServerID
	}
	ms, err := p.smsCore.GetMediaServer(context.Background(), mediaID)
	if err != nil {
		slog.Warn("onvif build rtsp: media server", "err", err, "channel", ch.ID)
		return ""
	}
	host := rtspHost(ms, p.conf)
	port := ms.Ports.RTSP
	if port == 0 {
		port = 554
	}
	return fmt.Sprintf("rtsp://%s:%d/%s/%s", host, port, ch.GetApp(), ch.GetStream())
}

// rtspHost 选取客户端可访问的流媒体 IP，避免返回本机回环地址。
func rtspHost(ms *sms.MediaServer, cfg *conf.Bootstrap) string {
	if ms.StreamIP != "" && isLANHost(ms.StreamIP) {
		return ms.StreamIP
	}
	if isLANHost(ms.IP) {
		return ms.IP
	}
	if h := advertiseHost(cfg); h != "" {
		return h
	}
	if ms.IP != "" {
		return ms.IP
	}
	return "127.0.0.1"
}
