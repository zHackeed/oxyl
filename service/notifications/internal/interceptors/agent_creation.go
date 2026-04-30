package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/notifications/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.AgentCreation] = (*AgentCreationInterceptor)(nil)

type AgentCreationInterceptor struct {
	agents *provider.AgentCompanyProvider
}

func NewAgentCreationInterceptor(agents *provider.AgentCompanyProvider) *AgentCreationInterceptor {
	return &AgentCreationInterceptor{agents: agents}
}

func (a *AgentCreationInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentCreation
}

func (a *AgentCreationInterceptor) Intercept(ctx context.Context, msg redisModels.AgentCreation) error {
	return a.agents.Add(ctx, msg.AgentId)
}
