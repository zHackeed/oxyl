package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.CompanyWebhookCreation] = (*WebhookCreation)(nil)

type WebhookCreation struct {
	settings *provider.NotificationSettingsProvider
}

func NewWebhookCreation(settings *provider.NotificationSettingsProvider) *WebhookCreation {
	return &WebhookCreation{settings: settings}
}

func (c *WebhookCreation) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyWebhookCreate
}

func (c *WebhookCreation) Intercept(ctx context.Context, msg redisModels.CompanyWebhookCreation) error {
	return c.settings.New(ctx, msg.Endpoint)
}
