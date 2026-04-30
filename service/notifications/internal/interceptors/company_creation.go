package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.CompanyCreation] = (*CompanyCreationInterceptor)(nil)

type CompanyCreationInterceptor struct {
	settings *provider.NotificationSettingsProvider
}

func NewCompanyCreationInterceptor(settings *provider.NotificationSettingsProvider) *CompanyCreationInterceptor {
	return &CompanyCreationInterceptor{settings: settings}
}

func (c *CompanyCreationInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyCreation
}

func (c *CompanyCreationInterceptor) Intercept(ctx context.Context, msg redisModels.CompanyCreation) error {
	return c.settings.Add(ctx, msg.CompanyId)
}
