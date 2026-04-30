package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"runtime"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
	"zhacked.me/oxyl/service/thresholds/internal/models"
	"zhacked.me/oxyl/service/thresholds/internal/provider"
	"zhacked.me/oxyl/service/thresholds/internal/storage"
	"zhacked.me/oxyl/shared/pkg/datasource"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	comm "zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

type MetricsAverageConsumer struct {
	agents     *provider.AgentMetadataProvider
	thresholds *provider.ThresholdProvider
	storage    *storage.ThresholdStorage
	redis      *datasource.RedisConnection
}

func NewMetricsConsumer(agents *provider.AgentMetadataProvider, thresholds *provider.ThresholdProvider, storage *storage.ThresholdStorage, redis *datasource.RedisConnection) *MetricsAverageConsumer {
	return &MetricsAverageConsumer{
		agents:     agents,
		thresholds: thresholds,
		storage:    storage,
		redis:      redis,
	}
}

func (m *MetricsAverageConsumer) Run(ctx context.Context) {
	workers := max(1, int(math.Floor(float64(runtime.NumCPU())*0.75)))

	for i := range workers {
		go m.worker(ctx, i, workers)
	}

	<-ctx.Done()
}

func (m *MetricsAverageConsumer) worker(ctx context.Context, index, total int) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			agentIDs := m.agents.IDs()
			for i := index; i < len(agentIDs); i += total {
				m.check(ctx, agentIDs[i])
			}
		}
	}
}

func (m *MetricsAverageConsumer) check(ctx context.Context, agentID string) {
	alive, err := m.redis.ExistAny(ctx, fmt.Sprintf(string(variables.RedisKeyHeartbeat), agentID))
	if err != nil || !alive {
		m.clearActiveThresholds(ctx, agentID)
		return
	}

	agentMetadata, found := m.agents.Get(agentID)
	if !found {
		return
	}

	companyThresholds, found := m.thresholds.Get(agentMetadata.Holder)
	if !found {
		return
	}

	var (
		generalMetrics    *models.GeneralMetricsAvg
		mountPointMetrics []*models.DiskUsageAvg
		networkMetrics    []*models.NetworkAvg
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		generalMetrics, err = m.storage.GetGeneralAvg(gctx, agentID, 30)
		return err
	})

	g.Go(func() error {
		var err error
		mountPointMetrics, err = m.storage.GetDiskUsageAvg(gctx, agentID, 30)
		return err
	})

	g.Go(func() error {
		var err error
		networkMetrics, err = m.storage.GetNetworkAvg(gctx, agentID, 30)
		return err
	})

	if err := g.Wait(); err != nil {
		slog.Error("unable to get metrics", "err", err)
		return
	}

	usedMemory := (generalMetrics.AvgMemory / float64(agentMetadata.TotalMemory*1024)) * 100

	m.evaluate(ctx, agentID, comm.NotificationTypeAgentCpuUsageThreshold,
		generalMetrics.AvgCPU,
		companyThresholds.ExceedsCPU(generalMetrics.AvgCPU),
	)

	m.evaluate(ctx, agentID, comm.NotificationTypeAgentMemoryUsageThreshold,
		usedMemory,
		companyThresholds.ExceedsMemory(usedMemory),
	)

	for _, metric := range mountPointMetrics {
		mountPointSize := agentMetadata.Partitions[metric.MountPoint]
		percentUsage := (metric.AvgDiskUsage / float64(mountPointSize.TotalSize)) * 100
		m.evaluate(ctx, agentID, comm.NotificationTypeAgentDiskUsageThreshold,
			percentUsage,
			companyThresholds.ExceedsDisk(percentUsage),
		)
	}

	for _, metric := range networkMetrics {
		rxMbps := (metric.AvgRXRate * 8) / 1_000_000
		txMbps := (metric.AvgTXRate * 8) / 1_000_000

		m.evaluate(ctx, agentID, comm.NotificationTypeAgentNetworkUsageThreshold,
			rxMbps,
			companyThresholds.ExceedsNetwork(rxMbps, txMbps),
		)
	}
}

func (m *MetricsAverageConsumer) evaluate(
	ctx context.Context,
	agentID string,
	reason comm.NotificationType,
	value float64,
	exceeded bool,
) error {
	key := thresholdActiveKey(agentID, reason)
	slog.Info("metric", "agentID", agentID, "reason", reason, "value", value, "exceeded", exceeded, "key", key)

	if exceeded {
		identifier := ulid.Make().String()

		hadValue, err := m.redis.SetNX(ctx, key, identifier, 2*time.Minute)
		if err != nil {
			return err
		}

		if hadValue != "" {
			m.redis.UpdateTTL(ctx, key, 2*time.Minute)
		}

		return m.redis.Publish(ctx, variables.RedisChannelThresholdNotification, redisModels.ThresholdNotification{
			Identifier:    identifier,
			AgentID:       agentID,
			TriggerReason: reason,
			TriggerValue:  fmt.Sprintf("%.2f", value),
			Resolved:      false,
		})
	}

	identifier, err := m.redis.GetAndDelete(ctx, key)

	if err != nil {
		return err
	}

	if identifier == "" {
		return nil
	}

	return m.redis.Publish(ctx, variables.RedisChannelThresholdNotification, redisModels.ThresholdNotification{
		Identifier:    identifier,
		AgentID:       agentID,
		TriggerReason: reason,
		Resolved:      true,
	})
}

func (m *MetricsAverageConsumer) clearActiveThresholds(ctx context.Context, agentID string) {
	keys := make([]string, len(comm.NotificationTypes()))

	for _, t := range comm.NotificationTypes() {
		keys = append(keys, thresholdActiveKey(agentID, t))
	}

	m.redis.DelAny(ctx, keys...)
}

func thresholdActiveKey(agentID string, t comm.NotificationType) string {
	return fmt.Sprintf(string(variables.RedisKeyThresholdActive), agentID, t)
}
