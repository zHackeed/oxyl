package interceptors

import (
	"context"

	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

// If we get this message, it means that a new agent has been enrolled and we need to add it to the provider
// As depend on it.

var _ messenger.Interceptor[models.AgentEnrollment] = (*AgentEnrollmentInterceptor)(nil)

type AgentEnrollmentInterceptor struct {
	agentMetadataProvider *provider.AgentMetadataProvider
}

func NewAgentEnrollmentInterceptor(agentMetadataProvider *provider.AgentMetadataProvider) *AgentEnrollmentInterceptor {
	return &AgentEnrollmentInterceptor{
		agentMetadataProvider: agentMetadataProvider,
	}
}

func (a *AgentEnrollmentInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentEnrollment
}

func (a *AgentEnrollmentInterceptor) Intercept(ctx context.Context, msg models.AgentEnrollment) error {
	return a.agentMetadataProvider.Add(ctx, msg.AgentId)
}
