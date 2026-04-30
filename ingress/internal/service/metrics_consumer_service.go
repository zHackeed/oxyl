package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	protocolV1 "zhacked.me/oxyl/protocol/v1"
	"zhacked.me/oxyl/protocol/v1/monitoring"
	"zhacked.me/oxyl/shared/pkg/datasource"
	redisMessage "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/utils"
	"zhacked.me/oxyl/shared/pkg/variables"
)

type MetricsConsumerService struct {
	// todo: metrics storage

	metricsStorage *storage.MonitoringStorage
	messenger      *datasource.RedisConnection

	activeListenerMu sync.RWMutex
	activeListeners  map[string]bool

	protocolV1.UnimplementedMonitoringServiceServer
}

var _ protocolV1.MonitoringServiceServer = (*MetricsConsumerService)(nil)

func NewMetricsConsumerService(metricStorage *storage.MonitoringStorage, redis *datasource.RedisConnection) *MetricsConsumerService {
	return &MetricsConsumerService{
		metricsStorage:  metricStorage,
		messenger:       redis,
		activeListeners: make(map[string]bool),
	}
}

func (m *MetricsConsumerService) AddListener(agentId string) {
	m.activeListenerMu.Lock()
	m.activeListeners[agentId] = true
	m.activeListenerMu.Unlock()
}

func (m *MetricsConsumerService) RemoveListener(agentId string) {
	m.activeListenerMu.Lock()
	delete(m.activeListeners, agentId)
	m.activeListenerMu.Unlock()
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

	m.activeListenerMu.RLock()
	_, listening := m.activeListeners[agentId]
	m.activeListenerMu.RUnlock()

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

	if listening {
		m.notifyListener(ctx, agentId, convertedGeneralData, convertedMountPointMetrics, convertedNetworkMetrics)
	}

	m.messenger.SetAny(ctx, fmt.Sprintf(string(variables.RedisKeyHeartbeat), agentId), "1", 2*time.Minute)

	//slog.Info("metrics registered", slog.String("agent_id", agentId))

	return &monitoring.AgentMetricsResponse{
		Success: true,
	}, nil
}

func (m *MetricsConsumerService) notifyListener(ctx context.Context, agentId string, generalMetrics *models.AgentGeneralMetrics, mountedMetrics []*models.AgentMountPointMetrics, networkMetrics []*models.AgentNetworkMetrics) {
	if err := m.messenger.Publish(ctx, variables.RedisChannelAgentMetrics, redisMessage.AgentMetricEntry{
		AgentId:        agentId,
		GeneralMetrics: generalMetrics,
		MountedMetrics: mountedMetrics,
		NetworkMetrics: networkMetrics,
	}); err != nil {
		slog.Error("failed to publish metrics", slog.String("agent_id", agentId), slog.String("error", err.Error()))
	}
}
