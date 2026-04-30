package interceptors

import (
	"context"
	"log/slog"
	"time"

	"zhacked.me/oxyl/service/notifications/internal/notificator"
	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/service/notifications/internal/storage"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.ThresholdNotification] = (*ThresholdNotificationInterceptor)(nil)

type ThresholdNotificationInterceptor struct {
	storage      *storage.NotificationStorage
	agents       *provider.AgentCompanyProvider
	agentStorage *storage.AgentToCompanyMapperStorage
	settings     *provider.NotificationSettingsProvider
}

func NewThresholdNotificationInterceptor(
	storage *storage.NotificationStorage,
	agents *provider.AgentCompanyProvider,
	agentStorage *storage.AgentToCompanyMapperStorage,
	settings *provider.NotificationSettingsProvider,
) *ThresholdNotificationInterceptor {
	return &ThresholdNotificationInterceptor{
		storage:      storage,
		agents:       agents,
		agentStorage: agentStorage,
		settings:     settings,
	}
}

func (t *ThresholdNotificationInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelThresholdNotification
}

func (t *ThresholdNotificationInterceptor) Intercept(ctx context.Context, msg redisModels.ThresholdNotification) error {
	companyID, ok := t.agents.Get(msg.AgentID)
	if !ok {
		slog.Warn("unknown agent in notification event", "agent", msg.AgentID)
		return nil
	}

	if msg.Resolved {
		if err := t.storage.Ack(ctx, msg.Identifier); err != nil {
			return err
		}
	} else {
		if err := t.storage.Insert(ctx, msg.Identifier, msg.AgentID, msg.TriggerReason, msg.TriggerValue); err != nil {
			return err
		}
	}

	settings, ok := t.settings.Get(companyID)
	if !ok {
		return nil
	}

	displayName, err := t.agentStorage.GetDisplayName(ctx, msg.AgentID)
	if err != nil {
		return err
	}

	for _, s := range settings {
		go t.deliverWithRetry(ctx, s.ID, func() error {
			return notificator.Send(s, msg, displayName)
		})
	}

	return nil
}

func (t *ThresholdNotificationInterceptor) deliverWithRetry(ctx context.Context, settingID string, fn func() error) {
	const maxAttempts = 5
	backoff := time.Second

	for range maxAttempts {
		if err := fn(); err == nil {
			return
		}

		slog.Warn("notification delivery failed, retrying", "setting", settingID)

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
			backoff *= 2
		}
	}

	slog.Error("notification delivery exhausted attempts", "setting", settingID)
}
