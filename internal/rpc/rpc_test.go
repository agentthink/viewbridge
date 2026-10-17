package rpc

import (
	"context"
	"log/slog"
	"testing"

	"github.com/gowvp/owl/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHealthCheck(t *testing.T) {
	addr := "localhost:50051"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	cli := protos.NewHealthClient(conn)
	resp, err := cli.Check(context.Background(), &protos.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("HealthCheck", "resp", resp)
}
