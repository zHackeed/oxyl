package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.CompanyWebhookDeletion] = (*WebhookDeletion)(nil)

type WebhookDeletion struct {
	settings *provider.NotificationSettingsProvider
}

func NewWebhookDeletion(settings *provider.NotificationSettingsProvider) *WebhookDeletion {
	return &WebhookDeletion{settings: settings}
}

func (c *WebhookDeletion) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyWebhookDelete
}

func (c *WebhookDeletion) Intercept(_ context.Context, msg redisModels.CompanyWebhookDeletion) error {
	c.settings.RemoveSetting(msg.CompanyId, msg.Endpoint)
	return nil
}
