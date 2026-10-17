package lalmax

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Config struct {
	URL    string
	Secret string
}

type Engine struct {
	cfg Config
	cli *http.Client
}

func NewEngine() Engine {
	return Engine{
		cli: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        30,
				MaxIdleConnsPerHost: 30,
				MaxConnsPerHost:     100,
			},
		},
	}
}

func (e Engine) SetConfig(cfg Config) Engine {
	e.cfg = cfg
	return e
}

// post 发送 POST 请求到 lalmax API
// 用法示例：e.post(ctx, "/api/path", map[string]any{"key": "value"}, &response)
func (e *Engine) post(ctx context.Context, path string, data map[string]any, out any) error {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, e.cfg.URL+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

// get 发送 GET 请求到 lalmax API
// 用法示例：e.get(ctx, "/api/path", &response)
func (e *Engine) get(ctx context.Context, path string, out any) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, e.cfg.URL+path, nil)
	resp, err := e.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}
