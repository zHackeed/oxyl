package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[models.CompanyDeletion] = (*CompanyDeletionInterceptor)(nil)

type CompanyDeletionInterceptor struct {
	companyMetadataProvider *provider.ThresholdProvider
}

func NewCompanyDeletionInterceptor(companyMetadataProvider *provider.ThresholdProvider) *CompanyDeletionInterceptor {
	return &CompanyDeletionInterceptor{
		companyMetadataProvider: companyMetadataProvider,
	}
}

func (c *CompanyDeletionInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyDeletion
}

func (c *CompanyDeletionInterceptor) Intercept(_ context.Context, msg models.CompanyDeletion) error {
	return c.companyMetadataProvider.Remove(msg.CompanyId)
}
