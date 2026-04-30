package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"zhacked.me/oxyl/service/thresholds/internal/consumer"
	"zhacked.me/oxyl/service/thresholds/internal/interceptors"
	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/service/thresholds/internal/storage"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/logger"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/version"
)

var thresholdConsumer = &cobra.Command{
	Use:   "",
	Short: "Threshold manager and notifier",
	Run:   startConsumingThresholds,
}

func init() {
	cobra.OnInitialize(func() {
		logger.Register(slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	})
}

func Execute(ctx context.Context) error {
	return thresholdConsumer.ExecuteContext(ctx)
}

func startConsumingThresholds(cmd *cobra.Command, args []string) {
	slog.Info("starting threshold manager server", slog.String("version", version.CommitID), slog.String("branch", version.Branch))

	initCtx, initCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer initCancel()

	timescale, redis, redisMessenger, err := createDatabases(initCtx)
	if err != nil {
		slog.Error("unable to connect to databases", err)
		return
	}
	defer timescale.Close()
	defer redis.Close()
	defer redisMessenger.Close()

	thresholdStorage, agentStorage, err := createStorages(timescale)
	if err != nil {
		slog.Error("unable to create threshold storage", err)
		return
	}

	thresholdProvider := provider.NewThresholdProvider(thresholdStorage)
	if err := thresholdProvider.Load(initCtx); err != nil {
		slog.Error("unable to load thresholds", err)
		return
	}

	agentProvider := provider.NewAgentProvider(agentStorage)
	if err := agentProvider.Load(initCtx); err != nil {
		slog.Error("unable to load agents", err)
		return
	}

	registerInterceptors(redisMessenger, thresholdProvider, agentProvider)

	go func() {
		if err := redisMessenger.Run(cmd.Context()); err != nil {
			slog.Error("unable to run messenger", err)
		}
	}()

	metricConsumer := consumer.NewMetricsConsumer(agentProvider, thresholdProvider, thresholdStorage, redis)
	metricConsumer.Run(cmd.Context())

	<-cmd.Context().Done()

	slog.Info("stopping threshold manager server")
}

func registerInterceptors(router *messenger.PubSubRouter, thresholdProvider *provider.ThresholdProvider, agentProvider *provider.AgentMetadataProvider) {
	messenger.RegisterHandler[models.AgentEnrollment](router, interceptors.NewAgentEnrollmentInterceptor(agentProvider))
	messenger.RegisterHandler[models.CompanyCreation](router, interceptors.NewCompanyCreationInterceptor(thresholdProvider))
	messenger.RegisterHandler[models.CompanyDeletion](router, interceptors.NewCompanyDeletionInterceptor(thresholdProvider))
	messenger.RegisterHandler[models.AgentDelete](router, interceptors.NewAgentRemovalInterceptor(agentProvider))
	messenger.RegisterHandler[models.ThresholdUpdate](router, interceptors.NewThresholdUpdateInterceptor(thresholdProvider))
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
		return nil, nil, nil, fmt.Errorf("unable to connect to redis for the messenger: %w", err)
	}

	return timescale, redis, redisMessenger, nil
}

func createStorages(timescale *datasource.TimescaleConnection) (*storage.ThresholdStorage, *storage.AgentStorage, error) {
	return storage.NewThresholdStorage(timescale), storage.NewAgentStorage(timescale), nil
}
