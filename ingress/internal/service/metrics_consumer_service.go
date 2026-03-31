package service

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	protocolV1 "zhacked.me/oxyl/protocol/v1"
	"zhacked.me/oxyl/protocol/v1/monitoring"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/utils"
)

type MetricsConsumerService struct {
	// todo: metrics storage
	protocolV1.UnimplementedMonitoringServiceServer
}

var _ protocolV1.MonitoringServiceServer = (*MetricsConsumerService)(nil)

func NewMetricsConsumerService() *MetricsConsumerService {
	return &MetricsConsumerService{}
}

func (m *MetricsConsumerService) ConsumeMetrics(ctx context.Context, in *monitoring.AgentMetrics) (*monitoring.AgentMetricsResponse, error) {
	agentId, found := utils.GetValueFromContext[string](ctx, models.ContextAgent)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "unauthenticated or missing context data")
	}

	slog.Info("metrics received", slog.String("agent_id", agentId), slog.Any("metrics", in))
	return &monitoring.AgentMetricsResponse{
		Success: true,
	}, nil
}
