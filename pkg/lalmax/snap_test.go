package lalmax

import (
	"context"
	"testing"
)

// TestGetSnapshot 测试获取快照图片
// 用法示例：go test -v -run TestGetSnapshot ./pkg/lalmax/
// 前提条件：
//   - lalmax 服务运行在 http://localhost:8080
//   - 有正在推流的流（如 live/test）
func TestGetSnapshot(t *testing.T) {
	ctx := context.Background()

	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "",
	})

	// 测试获取快照图片
	// 注意：需要有正在推流的流才能成功获取快照
	streamName := "live/test" // 替换为实际的流名称

	imageData, err := engine.GetSnapshot(ctx, streamName)
	if err != nil {
		// 如果流不存在，测试会失败但这是预期的
		t.Logf("GetSnapshot 失败（可能是流不存在）: %v", err)
		t.Skip("跳过测试：需要有正在推流的流")
		return
	}

	// 验证返回的数据
	if len(imageData) == 0 {
		t.Error("GetSnapshot 返回空数据")
		return
	}

	// 检查是否是 PNG 格式（PNG 文件头）
	if len(imageData) >= 8 {
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		isPNG := true
		for i := range 8 {
			if imageData[i] != pngHeader[i] {
				isPNG = false
				break
			}
		}
		if isPNG {
			t.Logf("GetSnapshot 成功，返回 PNG 图片，大小: %d bytes", len(imageData))
		} else {
			t.Logf("GetSnapshot 成功，返回数据大小: %d bytes（非 PNG 格式）", len(imageData))
		}
	}
}

// TestGetKeyFrameData 测试获取关键帧数据
// 用法示例：go test -v -run TestGetKeyFrameData ./pkg/lalmax/
func TestGetKeyFrameData(t *testing.T) {
	ctx := context.Background()

	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "",
	})

	streamName := "live/test" // 替换为实际的流名称

	keyframeData, err := engine.GetKeyFrameData(ctx, streamName)
	if err != nil {
		t.Logf("GetKeyFrameData 失败（可能是流不存在）: %v", err)
		t.Skip("跳过测试：需要有正在推流的流")
		return
	}

	if len(keyframeData) == 0 {
		t.Error("GetKeyFrameData 返回空数据")
		return
	}

	t.Logf("GetKeyFrameData 成功，返回数据大小: %d bytes", len(keyframeData))
}

// TestGetSnapshotInvalidStream 测试获取不存在的流的快照
func TestGetSnapshotInvalidStream(t *testing.T) {
	ctx := context.Background()

	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "",
	})

	// 使用一个不存在的流名称
	streamName := "nonexistent/stream_12345"

	_, err := engine.GetSnapshot(ctx, streamName)
	if err == nil {
		t.Error("期望返回错误，但没有")
		return
	}

	t.Logf("正确返回错误: %v", err)
}

// TestGetSnapshotEmptyStreamName 测试空流名称
func TestGetSnapshotEmptyStreamName(t *testing.T) {
	ctx := context.Background()

	engine := NewEngine()
	engine = engine.SetConfig(Config{
		URL:    "http://localhost:8080",
		Secret: "",
	})

	_, err := engine.GetSnapshot(ctx, "")
	if err == nil {
		t.Error("期望返回错误，但没有")
		return
	}

	t.Logf("正确返回错误: %v", err)
}
