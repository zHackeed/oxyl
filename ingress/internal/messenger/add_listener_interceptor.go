package messenger

import (
	"context"
	"log/slog"

	"zhacked.me/oxyl/ingress/internal/service"
	"zhacked.me/oxyl/shared/pkg/messenger"
	"zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

var _ messenger.Interceptor[models.AgentListening] = (*AddListenerInterceptor)(nil)

type AddListenerInterceptor struct {
	metricsConsumer *service.MetricsConsumerService
}

func NewAddListenerInterceptor(metricsConsumer *service.MetricsConsumerService) *AddListenerInterceptor {
	return &AddListenerInterceptor{
		metricsConsumer: metricsConsumer,
	}
}

func (a AddListenerInterceptor) GetChannel() variables.RedisChannel {
	return variables.RedisChannelAgentListening
}

func (a AddListenerInterceptor) Intercept(ctx context.Context, msg models.AgentListening) error {
	slog.Info("adding listener", slog.String("agent_id", msg.AgentId))
	a.metricsConsumer.AddListener(msg.AgentId)
	return nil
}
