package zlm

const (
	closeStreams = `/index/api/close_streams`
)

type CloseStreamsRequest struct {
	Schema string `json:"schema,omitempty"` // 协议名，如 rtsp/rtmp/hls，为空则关闭所有协议
	Vhost  string `json:"vhost,omitempty"`
	App    string `json:"app,omitempty"`
	Stream string `json:"stream,omitempty"`
	Force  bool   `json:"force,omitempty"` // 是否强制关闭（不等待注销超时）
}

type CloseStreamsResponse struct {
	Code          int `json:"code"`
	CountHit      int `json:"count_hit"`      // 筛选命中的流数量
	CountClosed   int `json:"count_closed"`    // 被关闭的流数量
}

// CloseStreams 关闭流
// https://docs.zlmediakit.com/zh/guide/media_server/restful_api.html#_7-index-api-close_streams
func (e *Engine) CloseStreams(in CloseStreamsRequest) (*CloseStreamsResponse, error) {
	body, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	var resp CloseStreamsResponse
	if err := e.post(closeStreams, body, &resp); err != nil {
		return nil, err
	}
	if err := e.ErrHandle(resp.Code, "close streams err"); err != nil {
		return nil, err
	}
	return &resp, nil
}
