package messenger

import (
	"context"
	"log/slog"

	"zhacked.me/oxyl/ingress/internal/service"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[models.AgentListening] = (*RemoveListenerInterceptor)(nil)

type RemoveListenerInterceptor struct {
	metricsConsumer *service.MetricsConsumerService
}

func NewRemoveListenerInterceptor(metricsConsumer *service.MetricsConsumerService) *RemoveListenerInterceptor {
	return &RemoveListenerInterceptor{
		metricsConsumer: metricsConsumer,
	}
}

func (r RemoveListenerInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentStoppedListening
}

func (r RemoveListenerInterceptor) Intercept(ctx context.Context, msg models.AgentListening) error {
	slog.Info("removing listener", slog.String("agent_id", msg.AgentId))
	r.metricsConsumer.RemoveListener(msg.AgentId)
	return nil
}
