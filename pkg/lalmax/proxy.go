package lalmax

import (
	"context"
	"encoding/json"
)

const (
	apiCtrlStartRelayPull = "/api/ctrl/start_relay_pull"
	apiCtrlStopRelayPull  = "/api/ctrl/stop_relay_pull"
)

type ApiCtrlStartRelayPullReq struct {
	Url           string `json:"url"`             // 必填项，回源拉流的完整url地址，目前支持rtmp和rtsp
	StreamName    string `json:"stream_name"`     // 选填项，如果不指定，则从`url`参数中解析获取
	PullTimeoutMs int    `json:"pull_timeout_ms"` // 选填项，pull建立会话的超时时间，单位毫秒。
	//. 选填项，pull连接失败或者中途断开连接的重试次数
	//  -1  表示一直重试，直到收到stop请求，或者开启并触发下面的自动关闭功能
	//  = 0 表示不重试
	//  > 0 表示重试次数
	// 提示：不开启自动重连，你可以在收到HTTP-Notify on_relay_pull_stop, on_update等消息时决定是否重连
	PullRetryNum int `json:"pull_retry_num"`
	// 选填项，没有观看者时，自动关闭pull会话，节约资源
	//  -1  表示不启动该功能
	//  = 0 表示没有观看者时，立即关闭pull会话
	//  > 0 表示没有观看者持续多长时间，关闭pull会话，单位毫秒
	//  默认值是-1
	//  提示：不开启该功能，你可以在收到HTTP-Notify on_sub_stop, on_update等消息时决定是否关闭relay pull
	AutoStopPullAfterNoOutMs int `json:"auto_stop_pull_after_no_out_ms"`
	//. 选填项，使用rtsp时的连接方式
	//  0 tcp
	//  1 udp
	//  默认值是0
	RtspMode int `json:"rtsp_mode"`
	//. 选填项，将接收的数据存成文件
	DebugDumpPacket          string `json:"debug_dump_packet"`
	KeepLiveFormGetParameter bool   `json:"keep_live_form_get_parameter"`
}

type CommonResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ApiCtrlStartRelayPullResp struct {
	CommonResp
	Data struct {
		StreamName string `json:"stream_name"`
		SessionId  string `json:"session_id"`
	} `json:"data"`
}

func (e *Engine) CtrlStartRelayPull(ctx context.Context, in ApiCtrlStartRelayPullReq) (*ApiCtrlStartRelayPullResp, error) {
	body, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	var resp ApiCtrlStartRelayPullResp
	if err := e.post(ctx, apiCtrlStartRelayPull, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func struct2map(in any) (map[string]any, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}
