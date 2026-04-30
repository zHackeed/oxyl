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
	agentListener "zhacked.me/oxyl/ingress/internal/messenger"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/logger"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
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

	timescale, redis, redisMessenger, err := createDatabases(initCtx)
	if err != nil {
		slog.Error("unable to create databases", "error", err)
		return
	}

	defer func() {
		if err := redis.Close(); err != nil {
			slog.Error("unable to close redis connection", "error", err)
		}

		if err := redisMessenger.Close(); err != nil {
			slog.Error("unable to close pubsub router", "error", err)
		}

		timescale.Close()
	}()

	companyStorage, agentStorage, notificationStorage, tokenStorage, monitoringStorage := createStorage(timescale, redis)
	agentService, tokenService, err := createServices(redis, companyStorage, agentStorage, notificationStorage, tokenStorage)

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
	metricsConsumerService := serviceV1.NewMetricsConsumerService(monitoringStorage, redis)

	messenger.RegisterHandler[models.AgentListening](
		redisMessenger, agentListener.NewAddListenerInterceptor(metricsConsumerService))

	messenger.RegisterHandler[models.AgentListening](
		redisMessenger, agentListener.NewRemoveListenerInterceptor(metricsConsumerService))

	go func() {
		if err := redisMessenger.Run(cmd.Context()); err != nil {
			slog.Error("unable to run pubsub router", "error", err)
		}
	}()

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
			interceptor.NewLoggingInterceptor().Intercept,
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

func createDatabases(ctx context.Context) (*datasource.TimescaleConnection, *datasource.RedisConnection, *messenger.PubSubRouter, error) {
	timescale, err := datasource.NewTimescaleConnection(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	redis, err := datasource.NewRedisConnection()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	redisMessenger, err := messenger.NewPubSubRouter()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	return timescale, redis, redisMessenger, nil
}

func createStorage(timescale *datasource.TimescaleConnection, redis *datasource.RedisConnection) (
	*storage.CompanyStorage, *storage.AgentStorage, *storage.NotificationStorage, *storage.TokenStorage, *storage.MonitoringStorage,
) {
	return storage.NewCompanyStorage(timescale),
		storage.NewAgentStorage(timescale),
		storage.NewNotificationStorage(timescale),
		storage.NewTokenStorage(redis),
		storage.NewMonitoringStorage(timescale)
}

func createServices(messenger *datasource.RedisConnection, companyStorage *storage.CompanyStorage, agentStorage *storage.AgentStorage, notificationStorage *storage.NotificationStorage, tokenStorage *storage.TokenStorage) (
	*service.AgentService, *service.TokenService, error,
) {
	agentService := service.NewAgentService(messenger, companyStorage, agentStorage, notificationStorage)

	tokenService, err := service.NewTokenService(tokenStorage, messenger)
	if err != nil {
		return nil, nil, err
	}

	return agentService, tokenService, nil
}
