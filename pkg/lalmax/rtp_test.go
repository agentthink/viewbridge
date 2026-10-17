package lalmax

import (
	"context"
	"testing"
)

// TestApiCtrlStartRtpPub 测试启动 RTP 发布
// 用法示例：go test -v -run TestApiCtrlStartRtpPub ./pkg/lalmax/
// 前提条件：需要 lalmax 服务运行在 http://localhost:8080
func TestApiCtrlStartRtpPub(t *testing.T) {
	ctx := context.Background()

	// 创建 Engine 实例
	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "", // 如果需要密钥，请填写
	})

	// 构造请求参数
	// 注意：实际使用中可能需要确保流名称对应的流存在或符合预期，
	// 这里仅测试接口调用的连通性和基本参数传递。
	req := ApiCtrlStartRtpPubReq{
		StreamName:      "test110", // 测试流名称
		Port:            0,         // 0 表示让服务器自动分配端口
		PeerPort:        0,         // 如不涉及对端被动接收，可填0或按需填写
		TimeoutMs:       10000,     // 超时时间
		IsTcpFlag:       0,         // 0: UDP, 1: TCP
		IsWaitKeyFrame:  1,         // 是否等待关键帧
		DebugDumpPacket: "",        // 调试抓包文件路径
		IsTcpActive:     false,     // 是否为 TCP 主动模式
	}

	// 调用 ApiCtrlStartRtpPub 方法
	resp, err := engine.ApiCtrlStartRtpPub(ctx, req)
	if err != nil {
		t.Fatalf("ApiCtrlStartRtpPub 调用失败: %v", err)
	}

	// 验证返回结果不为空
	if resp == nil {
		t.Fatal("ApiCtrlStartRtpPub 返回 nil")
	}

	// 打印返回结果，用于调试
	t.Logf("响应代码 (Code): %d", resp.Code)
	t.Logf("响应描述 (Msg): %s", resp.Msg)
	t.Logf("流名称 (StreamName): %s", resp.Data.StreamName)
	t.Logf("会话 ID (SessionId): %s", resp.Data.SessionId)
	t.Logf("端口 (Port): %d", resp.Data.Port)

	// 简单的校验，确保请求的流名称和返回的一致（如果成功的话）
	if resp.Code == 0 {
		if resp.Data.StreamName != req.StreamName {
			t.Errorf("期望流名称 %s, 实际返回 %s", req.StreamName, resp.Data.StreamName)
		}
		if resp.Data.Port == 0 && req.Port != 0 {
			// 如果请求指定了端口，期望返回该端口（视具体 lal 实现而定，这里仅作一般性检查）
			t.Logf("注意：返回端口为 0")
		}
	} else {
		t.Logf("注意：接口返回了非 0 错误码，这在没有实际流的情况下可能是预期的。")
	}
}
