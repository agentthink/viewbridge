package recording

import (
	"context"
	"log/slog"
	"strings"
)

// StartRecording 启动录制，在流注册时调用
// 根据配置决定是否录制该流，并通知 ZLM 开始 MP4 录制
func (c Core) StartRecording(ctx context.Context, channelType, app, stream string) error {
	if !c.IsEnabled() {
		slog.DebugContext(ctx, "录制未启用", "app", app, "stream", stream)
		return nil
	}

	if c.smsProvider == nil {
		slog.WarnContext(ctx, "SMS provider 未配置，无法启动录制")
		return nil
	}

	// 构建自定义存储路径：直接使用 storageDir
	// ZLM 会在此基础上创建 record/{app}/{stream}/{date}/ 目录结构
	customPath := c.conf.StorageDir
	maxSecond := c.conf.SegmentSeconds
	// 限制切片时长在 60 秒到 3600 秒之间
	maxSecond = max(maxSecond, 60)
	maxSecond = min(maxSecond, 3600)

	if err := c.smsProvider.StartRecord(app, stream, customPath, maxSecond); err != nil {
		slog.ErrorContext(ctx, "启动录制失败", "app", app, "stream", stream, "err", err)
		return err
	}

	slog.InfoContext(ctx, "启动录制成功", "app", app, "stream", stream, "path", customPath)
	return nil
}

// StopRecording 停止录制，在流注销时调用
// 流不存在时视为已停止（幂等），打 debug 而非 error
func (c Core) StopRecording(ctx context.Context, app, stream string) error {
	if c.smsProvider == nil {
		return nil
	}

	if err := c.smsProvider.StopRecord(app, stream); err != nil {
		// ZLM 的 stopRecord 接口在流不存在时返回此错误，属于正常的幂等语义
		if strings.Contains(err.Error(), "can not find the stream") {
			slog.DebugContext(ctx, "流已不存在，无需停止录制", "app", app, "stream", stream)
			return nil
		}
		slog.ErrorContext(ctx, "停止录制失败", "app", app, "stream", stream, "err", err)
		return err
	}

	slog.InfoContext(ctx, "停止录制成功", "app", app, "stream", stream)
	return nil
}
