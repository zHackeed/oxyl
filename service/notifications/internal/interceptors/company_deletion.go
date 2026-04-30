package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.CompanyDeletion] = (*CompanyDeletionInterceptor)(nil)

type CompanyDeletionInterceptor struct {
	settings *provider.NotificationSettingsProvider
}

func NewCompanyDeletionInterceptor(settings *provider.NotificationSettingsProvider) *CompanyDeletionInterceptor {
	return &CompanyDeletionInterceptor{settings: settings}
}

func (c *CompanyDeletionInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyDeletion
}

func (c *CompanyDeletionInterceptor) Intercept(_ context.Context, msg redisModels.CompanyDeletion) error {
	c.settings.Remove(msg.CompanyId)
	return nil
}
