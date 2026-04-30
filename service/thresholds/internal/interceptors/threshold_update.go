package interceptors

import (
	"context"
	"fmt"

	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.ThresholdUpdate] = (*ThresholdUpdateInterceptor)(nil)

type ThresholdUpdateInterceptor struct {
	companyMetadataProvider *provider.ThresholdProvider
}

func NewThresholdUpdateInterceptor(companyMetadataProvider *provider.ThresholdProvider) *ThresholdUpdateInterceptor {
	return &ThresholdUpdateInterceptor{
		companyMetadataProvider: companyMetadataProvider,
	}
}

func (t ThresholdUpdateInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelCompanyThresholdUpdate
}

func (t ThresholdUpdateInterceptor) Intercept(ctx context.Context, msg redisModels.ThresholdUpdate) error {
	companyMetadata, found := t.companyMetadataProvider.Get(msg.CompanyId)
	if !found {
		return fmt.Errorf("unable to find company metadata for company %q", msg.CompanyId)
	}

	switch msg.ThresholdType {
	case models.NotificationTypeAgentCpuUsageThreshold:
		companyMetadata.SetCPU(msg.Threshold)
	case models.NotificationTypeAgentMemoryUsageThreshold:
		companyMetadata.SetMemory(msg.Threshold)
	case models.NotificationTypeAgentDiskUsageThreshold:
		companyMetadata.SetMount(msg.Threshold)
	case models.NotificationTypeAgentDiskHealthThreshold:
		companyMetadata.SetDisk(msg.Threshold)
	case models.NotificationTypeAgentNetworkUsageThreshold:
		companyMetadata.SetNetworkTX(msg.Threshold)
		companyMetadata.SetNetworkRX(msg.Threshold)
	default:
		return fmt.Errorf("unknown threshold type %q", msg.ThresholdType)
	}

	return nil
}
