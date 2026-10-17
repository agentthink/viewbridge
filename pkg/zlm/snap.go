package zlm

import (
	"encoding/json"
	"fmt"

	"github.com/ixugo/goddd/pkg/hook"
)

const (
	getSnapshot = `/index/api/getSnap`
)

type GetSnapRequest struct {
	URL        string `json:"url"`
	TimeoutSec int    `json:"timeout_sec"` // 截图失败超时时间
	ExpireSec  int    `json:"expire_sec"`  // 截图过期时间
}

// GetSnap 获取截图或生成实时截图并返回
// https://docs.zlmediakit.com/zh/guide/media_server/restful_api.html#_23%E3%80%81-index-api-getsnap
func (e Engine) GetSnap(in GetSnapRequest) ([]byte, error) {
	body, err := struct2map(in)
	if err != nil {
		return nil, err
	}
	b, err := e.post2(getSnapshot, body)
	if err != nil {
		return nil, err
	}

	if len(b) < 100 {
		var resp OpenRTPServerResponse
		if err := json.Unmarshal(b, &resp); err == nil {
			if err := e.ErrHandle(resp.Code, resp.Msg); err != nil {
				return nil, err
			}
		}
	}
	if len(b) == 47255 && hook.MD5FromBytes(b) == "32ddfa5715059731ae893ec92fca0311" {
		return nil, fmt.Errorf("zlm: 没有更新图片")
	}
	return b, nil
}
