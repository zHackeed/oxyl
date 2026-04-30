package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[models.AgentDelete] = (*AgentRemovalInterceptor)(nil)

type AgentRemovalInterceptor struct {
	agentMetadataProvider *provider.AgentMetadataProvider
}

func NewAgentRemovalInterceptor(agentMetadataProvider *provider.AgentMetadataProvider) *AgentRemovalInterceptor {
	return &AgentRemovalInterceptor{
		agentMetadataProvider: agentMetadataProvider,
	}
}

func (a *AgentRemovalInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentDeletion
}

func (a *AgentRemovalInterceptor) Intercept(_ context.Context, msg models.AgentDelete) error {
	return a.agentMetadataProvider.Remove(msg.AgentId)
}
