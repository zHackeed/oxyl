package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"zhacked.me/oxyl/ingress/internal/interceptor"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/logger"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"

	serviceV1 "zhacked.me/oxyl/ingress/internal/service"
	protocolV1 "zhacked.me/oxyl/protocol/v1"
)

var nexus = cobra.Command{
	Use:   "",
	Short: "Start the ingress server",
	Run:   startNexus,
}

func init() {
	cobra.OnInitialize(func() {
		logger.Register(slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	})
}

func Execute(ctx context.Context) error {
	return nexus.ExecuteContext(ctx)
}

func startNexus(cmd *cobra.Command, _ []string) {
	slog.Info("starting ingress server")

	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	timescale, redis, err := createDatabases(initCtx)
	if err != nil {
		slog.Error("unable to create databases", "error", err)
		return
	}

	defer func() {
		if err := redis.Close(); err != nil {
			slog.Error("unable to close redis connection", "error", err)
		}

		timescale.Close()
	}()

	companyStorage, agentStorage, tokenStorage := createStorage(timescale, redis)
	agentService, tokenService, err := createServices(redis, companyStorage, agentStorage, tokenStorage)

	if err != nil {
		slog.Error("unable to create services", "error", err)
		return
	}

	agentAuthInterceptor := interceptor.NewAgentAuthInterceptor(tokenService, agentService)
	agentEnrollmentInterceptor, err := interceptor.NewAgentEnrollmentInterceptor(tokenService, agentService)
	if err != nil {
		slog.Error("unable to create agent enrollment interceptor", "error", err)
		return
	}

	defer agentEnrollmentInterceptor.Close()

	enrollmentService := serviceV1.NewEnrollmentService(redis, agentService, tokenService)
	metricsConsumerService := serviceV1.NewMetricsConsumerService()

	lis, err := net.Listen("tcp", ":19988")
	if err != nil {
		slog.Error("unable to listen on port 19988", "error", err)
		return
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			agentAuthInterceptor.Intercept,
			agentEnrollmentInterceptor.Intercept,
		),
	)

	protocolV1.RegisterEnrollmentServiceServer(grpcServer, enrollmentService)
	protocolV1.RegisterMonitoringServiceServer(grpcServer, metricsConsumerService)

	go func() {
		slog.Info("starting grpc server", slog.Int("port", 19988))
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("unable to start grpc server", "error", err)
		}
	}()

	<-cmd.Context().Done()

	grpcServer.GracefulStop()
	slog.Info("ingress server stopped")
}

func createDatabases(ctx context.Context) (*datasource.TimescaleConnection, *datasource.RedisConnection, error) {
	timescale, err := datasource.NewTimescaleConnection(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	redis, err := datasource.NewRedisConnection()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	return timescale, redis, nil
}

func createStorage(timescale *datasource.TimescaleConnection, redis *datasource.RedisConnection) (
	*storage.CompanyStorage, *storage.AgentStorage, *storage.TokenStorage,
) {
	return storage.NewCompanyStorage(timescale),
		storage.NewAgentStorage(timescale),
		storage.NewTokenStorage(redis)
}

func createServices(messenger *datasource.RedisConnection, companyStorage *storage.CompanyStorage, agentStorage *storage.AgentStorage, tokenStorage *storage.TokenStorage) (
	*service.AgentService, *service.TokenService, error,
) {
	agentService := service.NewAgentService(messenger, companyStorage, agentStorage)

	tokenService, err := service.NewTokenService(tokenStorage)
	if err != nil {
		return nil, nil, err
	}

	return agentService, tokenService, nil
}
