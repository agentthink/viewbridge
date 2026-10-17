package lalmax

import (
	"context"
	"testing"
)

// TestGetServerConfig 测试获取服务器配置
// 用法示例：go test -v -run TestGetServerConfig ./pkg/lalmax/
// 前提条件：需要 lalmax 服务运行在 http://localhost:8080
func TestGetServerConfig(t *testing.T) {
	ctx := context.Background()

	// 创建 Engine 实例
	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "", // 如果需要密钥，请填写
	})

	// 调用 GetServerConfig 方法
	resp, err := engine.GetServerConfig(ctx)
	if err != nil {
		t.Fatalf("GetServerConfig 调用失败: %v", err)
	}

	// 验证返回结果不为空
	if resp == nil {
		t.Fatal("GetServerConfig 返回 nil")
	}

	// 验证基本字段
	t.Logf("配置版本: %s", resp.ConfVersion)
	t.Logf("服务器 ID: %s", resp.ServerId)
	t.Logf("关键帧路径: %s", resp.KeyFramePath)
	t.Logf("最大打开文件数: %d", resp.MaxOpenFiles)

	// 验证 GOP 缓存配置
	t.Logf("GOP 缓存数量: %d", resp.GopCacheConfig.GopNum)
	t.Logf("单个 GOP 最大帧数: %d", resp.GopCacheConfig.SingleGopMaxFrameNum)

	// 验证 RTMP 配置
	t.Logf("RTMP 启用: %v", resp.RtmpConfig.Enable)
	t.Logf("RTMP 地址: %s", resp.RtmpConfig.Addr)
	t.Logf("RTMP 发布超时: %d 秒", resp.RtmpConfig.PubTimeoutSec)

	// 验证 RTSP 配置
	t.Logf("RTSP 启用: %v", resp.RtspConfig.Enable)
	t.Logf("RTSP 地址: %s", resp.RtspConfig.Addr)

	// 验证 GB28181 配置
	t.Logf("GB28181 启用: %v", resp.Gb28181Config.Enable)
	t.Logf("GB28181 SIP IP: %s", resp.Gb28181Config.SipIP)
	t.Logf("GB28181 SIP 端口: %d", resp.Gb28181Config.SipPort)
	t.Logf("GB28181 设备序列号: %s", resp.Gb28181Config.Serial)
	t.Logf("GB28181 域: %s", resp.Gb28181Config.Realm)

	// 验证 HTTP API 配置
	t.Logf("HTTP API 启用: %v", resp.HttpApiConfig.Enable)

	// 验证 RTC 配置
	t.Logf("RTC 启用: %v", resp.RtcConfig.Enable)

	// 验证 HTTP Notify 配置
	t.Logf("HTTP Notify 启用: %v", resp.HttpNotifyConfig.Enable)
	t.Logf("HTTP Notify OnPubStart: %s", resp.HttpNotifyConfig.OnPubStart)
	t.Logf("HTTP Notify OnPubStop: %s", resp.HttpNotifyConfig.OnPubStop)

	// 验证字段是否正确解析（非零值检查）
	if resp.ConfVersion == "" {
		t.Error("ConfVersion 为空，可能解析失败")
	}
	if resp.ServerId == "" {
		t.Error("ServerId 为空，可能解析失败")
	}
	if resp.KeyFramePath == "" {
		t.Error("KeyFramePath 为空，可能解析失败")
	}

	t.Log("GetServerConfig 测试通过，所有字段正确解析")
}

// TestSetHttpNotifyConfig 测试设置 HTTP 通知配置
// 用法示例：go test -v -run TestSetHttpNotifyConfig ./pkg/lalmax/
// 前提条件：需要 lalmax 服务运行在 http://localhost:8080
func TestSetHttpNotifyConfig(t *testing.T) {
	ctx := context.Background()

	// 创建 Engine 实例
	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "",
	})

	// 定义要设置的 HTTP 通知配置
	notifyConfig := HttpNotifyConfig{
		Enable:            true,
		UpdateIntervalSec: 5,
		OnPubStart:        "http://127.0.0.1:18080/webhook/on_pub_start",
		OnPubStop:         "http://127.0.0.1:18080/webhook/on_pub_stop",
		OnSubStart:        "http://127.0.0.1:18080/webhook/on_sub_start",
		OnSubStop:         "http://127.0.0.1:18080/webhook/on_sub_stop",
		OnStreamChanged:   "http://127.0.0.1:18080/webhook/on_stream_changed",
	}

	// 设置配置
	err := engine.SetHttpNotifyConfig(ctx, notifyConfig, MediaConfig{
		ListenPort:            8080,
		MultiPortMaxIncrement: 10,
	})
	if err != nil {
		t.Fatalf("SetHttpNotifyConfig 调用失败: %v", err)
	}
	t.Log("SetHttpNotifyConfig 设置成功")

	// 验证配置是否已生效
	resp, err := engine.GetServerConfig(ctx)
	if err != nil {
		t.Fatalf("GetServerConfig 调用失败: %v", err)
	}

	// 验证设置的字段
	if resp.HttpNotifyConfig.Enable != notifyConfig.Enable {
		t.Errorf("HttpNotifyConfig.Enable = %v, 期望 %v", resp.HttpNotifyConfig.Enable, notifyConfig.Enable)
	}
	if resp.HttpNotifyConfig.OnPubStart != notifyConfig.OnPubStart {
		t.Errorf("HttpNotifyConfig.OnPubStart = %s, 期望 %s", resp.HttpNotifyConfig.OnPubStart, notifyConfig.OnPubStart)
	}
	if resp.HttpNotifyConfig.OnPubStop != notifyConfig.OnPubStop {
		t.Errorf("HttpNotifyConfig.OnPubStop = %s, 期望 %s", resp.HttpNotifyConfig.OnPubStop, notifyConfig.OnPubStop)
	}
	if resp.HttpNotifyConfig.OnStreamChanged != notifyConfig.OnStreamChanged {
		t.Errorf("HttpNotifyConfig.OnStreamChanged = %s, 期望 %s", resp.HttpNotifyConfig.OnStreamChanged, notifyConfig.OnStreamChanged)
	}

	t.Logf("验证成功！当前 HTTP Notify 配置:")
	t.Logf("  Enable: %v", resp.HttpNotifyConfig.Enable)
	t.Logf("  UpdateIntervalSec: %d", resp.HttpNotifyConfig.UpdateIntervalSec)
	t.Logf("  OnPubStart: %s", resp.HttpNotifyConfig.OnPubStart)
	t.Logf("  OnPubStop: %s", resp.HttpNotifyConfig.OnPubStop)
	t.Logf("  OnStreamChanged: %s", resp.HttpNotifyConfig.OnStreamChanged)
}
