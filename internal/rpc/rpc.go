package rpc

import (
	"context"
	"log/slog"

	"github.com/gowvp/owl/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ protos.AnalysisServiceClient = (*AIClient)(nil)

// AIClient 封装 gRPC 检测服务客户端，提供统一的 AI 检测调用入口
type AIClient struct {
	cli protos.AnalysisServiceClient
}

// NewAIClient 创建 AI 检测客户端实例
func NewAIClient(addr string) *AIClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("NewAiClient", "err", err)
		return nil
	}

	go func() {
		p := protos.NewHealthClient(conn)
		resp, err := p.Check(context.Background(), &protos.HealthCheckRequest{})
		if err != nil {
			slog.Error("HealthCheck", "err", err)
			return
		}
		if resp.GetStatus() == protos.HealthCheckResponse_SERVING {
			slog.Info("HealthCheck OK", "resp", resp)
		} else {
			slog.Error("HealthCheck", "resp", resp)
		}
	}()

	cli := protos.NewAnalysisServiceClient(conn)
	return &AIClient{cli: cli}
}

// GetStatus implements [protos.AnalysisServiceClient].
func (a *AIClient) GetStatus(ctx context.Context, in *protos.StatusRequest, opts ...grpc.CallOption) (*protos.StatusResponse, error) {
	return a.cli.GetStatus(ctx, in, opts...)
}

// StartCamera implements [protos.AnalysisServiceClient].
func (a *AIClient) StartCamera(ctx context.Context, in *protos.StartCameraRequest, opts ...grpc.CallOption) (*protos.StartCameraResponse, error) {
	if in.GetDetectIntervalSeconds() <= 0 {
		in.DetectIntervalSeconds = 5.0
	}
	if in.GetAlertCooldownSeconds() <= 0 {
		in.AlertCooldownSeconds = 30.0
	}
	if in.GetThreshold() == 0 {
		in.Threshold = 0.5
	}
	if in.GetRetryLimit() == 0 {
		in.RetryLimit = 10
	}
	return a.cli.StartCamera(ctx, in, opts...)
}

// StopCamera implements [protos.AnalysisServiceClient].
func (a *AIClient) StopCamera(ctx context.Context, in *protos.StopCameraRequest, opts ...grpc.CallOption) (*protos.StopCameraResponse, error) {
	return a.cli.StopCamera(ctx, in, opts...)
}
