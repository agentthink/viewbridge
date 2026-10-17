package adapter

import (
	"github.com/gowvp/owl/internal/core/recording"
	"github.com/gowvp/owl/internal/core/sms"
	"github.com/gowvp/owl/pkg/zlm"
)

var _ recording.SMSProvider = (*SMSAdapter)(nil)

// SMSAdapter 实现 recording.SMSProvider 接口
// 将 sms.Core 的录制能力适配给 recording 领域使用
type SMSAdapter struct {
	smsCore sms.Core
}

// NewSMSAdapter 创建 SMS 适配器，返回 recording.SMSProvider 接口
// Wire 通过此函数自动绑定 sms.Core -> recording.SMSProvider
func NewSMSAdapter(smsCore sms.Core) recording.SMSProvider {
	return &SMSAdapter{smsCore: smsCore}
}

// StartRecord 启动录制
func (a *SMSAdapter) StartRecord(app, stream, customPath string, maxSecond int) error {
	ms, err := a.smsCore.GetDefaultMediaServer()
	if err != nil {
		return err
	}
	_, err = a.smsCore.StartRecord(ms, zlm.StartRecordRequest{
		Type:       1, // MP4
		Vhost:      "__defaultVhost__",
		App:        app,
		Stream:     stream,
		CustomPath: customPath,
		MaxSecond:  maxSecond,
	})
	return err
}

// StopRecord 停止录制
func (a *SMSAdapter) StopRecord(app, stream string) error {
	ms, err := a.smsCore.GetDefaultMediaServer()
	if err != nil {
		return err
	}
	_, err = a.smsCore.StopRecord(ms, zlm.StopRecordRequest{
		Type:   1, // MP4
		Vhost:  "__defaultVhost__",
		App:    app,
		Stream: stream,
	})
	return err
}

// ListRecordingStreams 批量获取所有在线流的录制状态
// 调用 ZLM getMediaList 一次获取全部流，提取 isRecordingMP4 状态
// 返回 map key 格式为 "app/stream"
func (a *SMSAdapter) ListRecordingStreams() (map[string]bool, error) {
	ms, err := a.smsCore.GetDefaultMediaServer()
	if err != nil {
		return nil, err
	}
	resp, err := a.smsCore.GetMediaList(ms)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(resp.Data))
	for _, item := range resp.Data {
		key := item.App + "/" + item.Stream
		if item.IsRecordingMP4 {
			result[key] = true
		} else if _, exists := result[key]; !exists {
			result[key] = false
		}
	}
	return result, nil
}
