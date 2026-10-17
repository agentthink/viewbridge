package zlm

const (
	openRtpServer  = `/index/api/openRtpServer`
	closeRtpServer = `/index/api/closeRtpServer`
)

type OpenRTPServerResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Port int    `json:"port"` // 接收端口，方便获取随机端口号
}
type OpenRTPServerRequest struct {
	Port     int    `json:"port"`      // 接收端口，0 则为随机端口
	TCPMode  int8   `json:"tcp_mode"`  // 0 udp 模式，1 tcp 被动模式, 2 tcp 主动模式。 (兼容 enable_tcp 为 0/1)
	StreamID string `json:"stream_id"` // 该端口绑定的流 ID，该端口只能创建这一个流(而不是根据 ssrc 创建多个)
	SSRC     uint64 `json:"ssrc"`      // 防串流：非零时 ZLM 只接收匹配该 SSRC 的 RTP 包，0 则不过滤
}

// OpenRTPServer 创建 GB28181 RTP 接收端口，如果该端口接收数据超时，则会自动被回收(不用调用 closeRtpServer 接口)
// https://docs.zlmediakit.com/zh/guide/media_server/restful_api.html#_24%E3%80%81-index-api-openrtpserver
func (e *Engine) OpenRTPServer(in OpenRTPServerRequest) (*OpenRTPServerResponse, error) {
	body, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	var resp OpenRTPServerResponse
	if err := e.post(openRtpServer, body, &resp); err != nil {
		return nil, err
	}
	if err := e.ErrHandle(resp.Code, resp.Msg); err != nil {
		return nil, err
	}
	return &resp, nil
}

type CloseRTPServerRequest struct {
	StreamID string `json:"stream_id"` // 调用 openRtpServer 接口时提供的流 ID
}

type CloseRTPServerResponse struct {
	Code int `json:"code"`
	Hit  int `json:"hit"` // 是否找到记录并关闭
}

// CloseRTPServer 关闭 GB28181 RTP 接收端口
// https://docs.zlmediakit.com/zh/guide/media_server/restful_api.html#_25%E3%80%81-index-api-closertpserver
func (e *Engine) CloseRTPServer(in CloseRTPServerRequest) (*CloseRTPServerResponse, error) {
	body, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	var resp CloseRTPServerResponse
	if err := e.post(closeRtpServer, body, &resp); err != nil {
		return nil, err
	}
	if err := e.ErrHandle(resp.Code, "rtp close err"); err != nil {
		return nil, err
	}
	return &resp, nil
}
