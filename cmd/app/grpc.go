package app

import (
	"context"
	"errors"
	"fmt"
	realtimev1 "middleware/api/realtime/v1"
	grpcserver "middleware/internal/api/realtime"
	"net"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const grpcServerAddress = ":1720"

func GRPCConfig(ctx context.Context) error {
	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		return fmt.Errorf("falha ao escutar em %s: %w", grpcServerAddress, err)
	}

	secretKey := viper.GetString("ACCESS_TOKEN")
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcserver.UnaryAuthInterceptor(secretKey)),
		grpc.StreamInterceptor(grpcserver.StreamAuthInterceptor(secretKey)),
	)
	realtimeServer := grpcserver.NewRealtimeTagServer()
	realtimev1.RegisterRealtimeTagServiceServer(server, realtimeServer)
	grpcserver.RegisterRealtimeTagCatalogServiceServer(server, realtimeServer)
	serveError := make(chan error, 1)
	go func() { serveError <- server.Serve(listener) }()

	select {
	case err := <-serveError:
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	stopped := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		server.Stop()
	}
	return nil
}
