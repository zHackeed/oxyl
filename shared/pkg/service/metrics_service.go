package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/utils"
)

type MetricsService struct {
	companyStorage      *storage.CompanyStorage
	agentStorage        *storage.AgentStorage
	agentMetricsStorage *storage.MonitoringStorage
}

func NewAgentMetricsService(companyStorage *storage.CompanyStorage, agentStorage *storage.AgentStorage, agentMetricStorage *storage.MonitoringStorage) *MetricsService {
	return &MetricsService{
		companyStorage:      companyStorage,
		agentStorage:        agentStorage,
		agentMetricsStorage: agentMetricStorage,
	}
}
func (s *MetricsService) GetMetrics(ctx context.Context, agentId string, interval time.Duration) (*models.AgentMetrics, error) {
	user, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	companyHolder, err := s.agentStorage.GetHolderOfAgent(ctx, agentId)
	if err != nil {
		return nil, fmt.Errorf("unable to get agent holder: %w", err)
	}

	memberOf, err := s.companyStorage.GetCompanyMembership(ctx, user, *companyHolder)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			return nil, models.ErrPermissionDenied
		}
		return nil, fmt.Errorf("unable to check if user is member of company: %w", err)
	}

	if !models.HasPermission(memberOf.Permission, models.CompanyPermissionView) {
		return nil, models.ErrPermissionDenied
	}

	var (
		generalMetrics []*models.AgentGeneralMetrics
		mountMetrics   []*models.AgentMountPointMetrics
		diskMetrics    []*models.AgentPhysicalDiskMetrics
		networkMetrics []*models.AgentNetworkMetrics
	)

	var g errgroup.Group

	g.Go(func() (err error) {
		generalMetrics, err = s.agentMetricsStorage.GetGeneralMetrics(ctx, agentId, interval)
		return
	})
	g.Go(func() (err error) {
		mountMetrics, err = s.agentMetricsStorage.GetMountPointMetrics(ctx, agentId, interval)
		return
	})
	g.Go(func() (err error) {
		diskMetrics, err = s.agentMetricsStorage.GetPhysicalDiskMetrics(ctx, agentId, interval)
		return
	})
	g.Go(func() (err error) {
		networkMetrics, err = s.agentMetricsStorage.GetNetworkMetrics(ctx, agentId, interval)
		return
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("unable to get metrics: %w", err)
	}

	composite, err := models.NewAgentMetrics(agentId, generalMetrics, mountMetrics, diskMetrics, networkMetrics, interval)
	if err != nil {
		return nil, fmt.Errorf("unable to create agent metrics: %w", err)
	}

	return composite, nil
}
