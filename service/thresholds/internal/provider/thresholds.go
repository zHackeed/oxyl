package provider

import (
	"context"
	"fmt"
	"sync"

	"zhacked.me/oxyl/service/thresholds/internal/models"
	"zhacked.me/oxyl/service/thresholds/internal/storage"
)

type ThresholdProvider struct {
	mu         sync.RWMutex
	thresholds map[string]*models.CompanyThresholds

	storage *storage.ThresholdStorage
}

func NewThresholdProvider(thresholdStorage *storage.ThresholdStorage) *ThresholdProvider {
	return &ThresholdProvider{
		thresholds: make(map[string]*models.CompanyThresholds),
		storage:    thresholdStorage,
	}
}

func (p *ThresholdProvider) Load(ctx context.Context) error {
	thresholds, err := p.storage.GetAllThresholds(ctx)

	if err != nil {
		return err
	}

	p.thresholds = thresholds

	return nil
}

func (p *ThresholdProvider) Get(companyID string) (*models.CompanyThresholds, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	t, ok := p.thresholds[companyID]
	return t, ok
}

func (p *ThresholdProvider) Add(ctx context.Context, companyID string) error {
	thresholds, err := p.storage.GetThresholds(ctx, companyID)
	if err != nil {
		return fmt.Errorf("unable to get thresholds for company %q: %w", companyID, err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.thresholds[companyID] = thresholds
	return nil
}

func (p *ThresholdProvider) Remove(companyId string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.thresholds, companyId)
	return nil
}
