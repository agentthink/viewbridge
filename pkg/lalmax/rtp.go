package lalmax

import "context"

const ctrlStartRtpPub = "/api/ctrl/start_rtp_pub"

type ApiCtrlStartRtpPubReq struct {
	StreamName      string `json:"stream_name"`
	Port            int    `json:"port"`
	PeerPort        int    `json:"peer_port"` // 对端收流端口
	TimeoutMs       int    `json:"timeout_ms"`
	IsTcpFlag       int    `json:"is_tcp_flag"`
	IsWaitKeyFrame  int    `json:"is_wait_key_frame"`
	DebugDumpPacket string `json:"debug_dump_packet"`
	IsTcpActive     bool   `json:"is_tcp_active"` // Tcp主动模式
}

type ApiCtrlStartRtpPubResp struct {
	CommonResp
	Data struct {
		StreamName string `json:"stream_name"`
		SessionId  string `json:"session_id"`
		Port       int    `json:"port"`
	} `json:"data"`
}

func (e *Engine) ApiCtrlStartRtpPub(ctx context.Context, in ApiCtrlStartRtpPubReq) (*ApiCtrlStartRtpPubResp, error) {
	body, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	var resp ApiCtrlStartRtpPubResp
	if err := e.post(ctx, ctrlStartRtpPub, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
