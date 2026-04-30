package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[models.CompanyCreation] = (*CompanyCreationInterceptor)(nil)

type CompanyCreationInterceptor struct {
	companyMetadataProvider *provider.ThresholdProvider
}

func NewCompanyCreationInterceptor(companyMetadataProvider *provider.ThresholdProvider) *CompanyCreationInterceptor {
	return &CompanyCreationInterceptor{
		companyMetadataProvider: companyMetadataProvider,
	}
}

func (c *CompanyCreationInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyCreation
}

func (c *CompanyCreationInterceptor) Intercept(ctx context.Context, msg models.CompanyCreation) error {
	return c.companyMetadataProvider.Add(ctx, msg.CompanyId)
}
