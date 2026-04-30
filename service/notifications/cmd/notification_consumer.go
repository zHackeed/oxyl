package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"zhacked.me/oxyl/service/notifications/internal/interceptors"
	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/service/notifications/internal/storage"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/logger"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/version"
)

var notificationConsumer = &cobra.Command{
	Use:   "",
	Short: "Notification consumer and dispatcher",
	Run:   startNotificationConsumer,
}

func init() {
	cobra.OnInitialize(func() {
		logger.Register(slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	})
}

func Execute(ctx context.Context) error {
	return notificationConsumer.ExecuteContext(ctx)
}

func startNotificationConsumer(cmd *cobra.Command, args []string) {
	slog.Info("starting notification consumer", slog.String("version", version.CommitID), slog.String("branch", version.Branch))

	initCtx, initCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer initCancel()

	timescale, redis, redisMessenger, err := createDatabases(initCtx)
	if err != nil {
		slog.Error("unable to connect to databases", "err", err)
		return
	}
	defer timescale.Close()
	defer redis.Close()
	defer redisMessenger.Close()

	notificationStorage, notificationSettingsStorage, agentStorage := createStorages(timescale)

	agentProvider := provider.NewAgentCompanyProvider(agentStorage)
	settingsProvider := provider.NewNotificationSettingsProvider(notificationSettingsStorage)

	if err := agentProvider.Load(initCtx); err != nil {
		slog.Error("unable to load agents", "err", err)
		return
	}

	if err := settingsProvider.Load(initCtx); err != nil {
		slog.Error("unable to load notification settings", "err", err)
		return
	}

	registerInterceptors(redisMessenger, notificationStorage, agentProvider, agentStorage, settingsProvider)

	go func() {
		if err := redisMessenger.Run(cmd.Context()); err != nil {
			slog.Error("unable to run messenger", "err", err)
		}
	}()

	<-cmd.Context().Done()

	slog.Info("stopping notification consumer")
}

func registerInterceptors(
	router *messenger.PubSubRouter,
	notificationStorage *storage.NotificationStorage,
	agents *provider.AgentCompanyProvider,
	agentStorage *storage.AgentToCompanyMapperStorage,
	settings *provider.NotificationSettingsProvider,
) {
	messenger.RegisterHandler[redisModels.AgentCreation](router, interceptors.NewAgentCreationInterceptor(agents))
	messenger.RegisterHandler[redisModels.AgentDelete](router, interceptors.NewAgentRemovalInterceptor(agents))
	messenger.RegisterHandler[redisModels.CompanyCreation](router, interceptors.NewCompanyCreationInterceptor(settings))
	messenger.RegisterHandler[redisModels.CompanyDeletion](router, interceptors.NewCompanyDeletionInterceptor(settings))
	messenger.RegisterHandler[redisModels.ThresholdNotification](router, interceptors.NewThresholdNotificationInterceptor(notificationStorage, agents, agentStorage, settings))
	messenger.RegisterHandler[redisModels.CompanyWebhookCreation](router, interceptors.NewWebhookCreation(settings))
	messenger.RegisterHandler[redisModels.CompanyWebhookDeletion](router, interceptors.NewWebhookDeletion(settings))
}

func createDatabases(ctx context.Context) (*datasource.TimescaleConnection, *datasource.RedisConnection, *messenger.PubSubRouter, error) {
	timescale, err := datasource.NewTimescaleConnection(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to connect to timescale: %w", err)
	}

	redis, err := datasource.NewRedisConnection()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to connect to redis: %w", err)
	}

	redisMessenger, err := messenger.NewPubSubRouter()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to connect to redis messenger: %w", err)
	}

	return timescale, redis, redisMessenger, nil
}

func createStorages(timescale *datasource.TimescaleConnection) (
	*storage.NotificationStorage,
	*storage.NotificationSettingStorage,
	*storage.AgentToCompanyMapperStorage,
) {
	return storage.NewNotificationStorage(timescale),
		storage.NewNotificationSettingStorage(timescale),
		storage.NewAgentToCompanyMapperStorage(timescale)
}
