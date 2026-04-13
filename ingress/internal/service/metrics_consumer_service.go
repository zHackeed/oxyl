package service

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	protocolV1 "zhacked.me/oxyl/protocol/v1"
	"zhacked.me/oxyl/protocol/v1/monitoring"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/utils"
)

type MetricsConsumerService struct {
	// todo: metrics storage

	metricsStorage *storage.MonitoringStorage

	protocolV1.UnimplementedMonitoringServiceServer
}

var _ protocolV1.MonitoringServiceServer = (*MetricsConsumerService)(nil)

func NewMetricsConsumerService(metricStorage *storage.MonitoringStorage) *MetricsConsumerService {
	return &MetricsConsumerService{
		metricsStorage: metricStorage,
	}
}

func (m *MetricsConsumerService) SendMetrics(ctx context.Context, in *monitoring.AgentMetrics) (*monitoring.AgentMetricsResponse, error) {
	agentId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyAgent)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "unauthenticated or missing context data")
	}

	convertedGeneralData, err := models.NewGeneralMetrics(in.GeneralMetrics.CpuUsage, in.GeneralMetrics.MemoryUsage, in.GeneralMetrics.Uptime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("failed to convert general metrics: %s", err.Error()))
	}

	convertedMountPointMetrics := make([]*models.AgentMountPointMetrics, 0)
	convertedPhysicalDiskMetrics := make([]*models.AgentPhysicalDiskMetrics, 0)
	convertedNetworkMetrics := make([]*models.AgentNetworkMetrics, 0)

	for _, mountPoint := range in.DiskMetrics {
		mountPoint, err := models.NewMountPointMetrics(mountPoint.MountPoint, mountPoint.UsedSpace)
		if err != nil {
			continue
		}

		convertedMountPointMetrics = append(convertedMountPointMetrics, mountPoint)
	}

	for _, blockDevice := range in.PhysicalDiskMetrics {
		blockDeviceMetric, err := models.NewAgentPhysicalDiskMetrics(blockDevice.DiskPath, blockDevice.HealthUsed,
			blockDevice.MediaErrors_1, blockDevice.MediaErrors_2,
			blockDevice.ErrorRate, blockDevice.PendingSectors)

		if err != nil {
			continue
		}

		convertedPhysicalDiskMetrics = append(convertedPhysicalDiskMetrics, blockDeviceMetric)
	}

	for _, ifData := range in.NetworkMetrics {
		ifData, err := models.NewAgentNetworkMetrics(ifData.InterfaceName,
			ifData.BytesReceived, ifData.BytesSent,
			ifData.PacketsReceived, ifData.PacketsSent,
			ifData.BytesReceivedRate, ifData.BytesSentRate,
			ifData.PacketsReceivedRate, ifData.PacketsSentRate)

		if err != nil {
			continue
		}

		convertedNetworkMetrics = append(convertedNetworkMetrics, ifData)
	}

	if err := m.metricsStorage.InsertData(ctx, agentId, convertedGeneralData, convertedMountPointMetrics, convertedPhysicalDiskMetrics, convertedNetworkMetrics); err != nil {
		slog.Error("failed to insert metrics", slog.String("agent_id", agentId), slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to insert metrics")
	}

	// todo: send response if they are listening on the socket

	slog.Info("metrics registered", slog.String("agent_id", agentId))

	return &monitoring.AgentMetricsResponse{
		Success: true,
	}, nil
}
