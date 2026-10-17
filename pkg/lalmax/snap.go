package lalmax

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	apiStatKeyFrame = "/api/stat/key_frame"
)

// SnapType 快照类型
type SnapType string

const (
	// SnapTypeImage 返回 PNG 图片
	SnapTypeImage SnapType = "image"
	// SnapTypeKeyFrame 返回关键帧数据（默认）
	SnapTypeKeyFrame SnapType = "keyframe"
)

// GetKeyFrameImage 获取流的关键帧图片（PNG 格式）
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	imageData, err := engine.GetKeyFrameImage(ctx, "live/test")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// imageData 是 PNG 图片的二进制数据
//	os.WriteFile("snapshot.png", imageData, 0644)
func (e *Engine) GetKeyFrameImage(ctx context.Context, streamName string) ([]byte, error) {
	return e.getKeyFrame(ctx, streamName, SnapTypeImage)
}

// GetKeyFrameData 获取流的关键帧数据
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	keyframeData, err := engine.GetKeyFrameData(ctx, "live/test")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (e *Engine) GetKeyFrameData(ctx context.Context, streamName string) ([]byte, error) {
	return e.getKeyFrame(ctx, streamName, SnapTypeKeyFrame)
}

// GetSnapshot 获取流的快照图片（GetKeyFrameImage 的别名，更直观的命名）
// 用法示例：
//
//	engine := lalmax.NewEngine().SetConfig(lalmax.Config{URL: "http://localhost:8080"})
//	imageData, err := engine.GetSnapshot(ctx, "live/test")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (e *Engine) GetSnapshot(ctx context.Context, streamName string) ([]byte, error) {
	return e.GetKeyFrameImage(ctx, streamName)
}

// getKeyFrame 内部方法：获取关键帧数据或图片
func (e *Engine) getKeyFrame(ctx context.Context, streamName string, snapType SnapType) ([]byte, error) {
	if streamName == "" {
		return nil, fmt.Errorf("lalmax: stream_name is required")
	}

	// 构建请求 URL
	url := fmt.Sprintf("%s%s?stream_name=%s", e.cfg.URL, apiStatKeyFrame, streamName)
	if snapType == SnapTypeImage {
		url += "&type=image"
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("lalmax: create request failed: %w", err)
	}

	// 发送请求
	resp, err := e.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lalmax: request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lalmax: read response failed: %w", err)
	}

	// 处理不同的 HTTP 状态码
	switch resp.StatusCode {
	case http.StatusOK:
		// 检查是否是 JSON 错误响应（服务端可能在 200 状态码下返回 JSON 错误）
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") || isJSONResponse(body) {
			var errResp FixedHeader
			if err := json.Unmarshal(body, &errResp); err == nil && errResp.Code != 0 && errResp.Code != 10000 {
				return nil, fmt.Errorf("lalmax: %s", errResp.Msg)
			}
		}
		return body, nil
	case http.StatusTooManyRequests: // 429
		return nil, fmt.Errorf("lalmax: keyframe is being generated, please try again later")
	case http.StatusNotFound:
		return nil, fmt.Errorf("lalmax: stream not found: %s", streamName)
	default:
		return nil, fmt.Errorf("lalmax: unexpected status code %d: %s", resp.StatusCode, string(body))
	}
}

// isJSONResponse 检查响应体是否是 JSON 格式
func isJSONResponse(body []byte) bool {
	if len(body) == 0 {
		return false
	}
	// JSON 通常以 { 或 [ 开头
	return body[0] == '{' || body[0] == '['
}
