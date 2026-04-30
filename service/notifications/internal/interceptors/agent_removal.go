package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.AgentDelete] = (*AgentRemovalInterceptor)(nil)

type AgentRemovalInterceptor struct {
	agents *provider.AgentCompanyProvider
}

func NewAgentRemovalInterceptor(agents *provider.AgentCompanyProvider) *AgentRemovalInterceptor {
	return &AgentRemovalInterceptor{agents: agents}
}

func (a *AgentRemovalInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentDeletion
}

func (a *AgentRemovalInterceptor) Intercept(_ context.Context, msg redisModels.AgentDelete) error {
	a.agents.Remove(msg.AgentId)
	return nil
}
