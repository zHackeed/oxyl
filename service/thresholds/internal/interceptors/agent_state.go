package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[redisModels.AgentStateUpdate] = (*AgentStateInterceptor)(nil)

type AgentStateInterceptor struct {
	agents *provider.AgentMetadataProvider
}

func NewAgentStateInterceptor(agents *provider.AgentMetadataProvider) *AgentStateInterceptor {
	return &AgentStateInterceptor{
		agents: agents,
	}
}

func (a AgentStateInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentStateUpdate
}

func (a AgentStateInterceptor) Intercept(ctx context.Context, msg redisModels.AgentStateUpdate) error {
	switch msg.Status {
	case models.AgentStatusActive:
		// The agent is active, we need to add it to the provider
		return a.agents.Add(ctx, msg.AgentId)
	case models.AgentStatusInactive:
	case models.AgentStatusMaintenance:
		// Do not track maintenance agents
		return a.agents.Remove(msg.AgentId)
	default:
		return nil
	}

	return nil
}
