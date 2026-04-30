package provider

import (
	"context"
	"fmt"
	"sync"

	"zhacked.me/oxyl/service/thresholds/internal/models"
	"zhacked.me/oxyl/service/thresholds/internal/storage"
)

type AgentMetadataProvider struct {
	mu     sync.RWMutex
	agents map[string]*models.AgentMetadata

	storage *storage.AgentStorage
}

func NewAgentProvider(agentMetadataStorage *storage.AgentStorage) *AgentMetadataProvider {
	return &AgentMetadataProvider{
		agents:  make(map[string]*models.AgentMetadata),
		storage: agentMetadataStorage,
	}
}

func (p *AgentMetadataProvider) Load(ctx context.Context) error {
	agents, err := p.storage.GetAll(ctx)

	if err != nil {
		return err
	}

	p.agents = agents

	return nil
}

func (p *AgentMetadataProvider) Get(id string) (*models.AgentMetadata, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	agentData, found := p.agents[id]
	return agentData, found
}

func (p *AgentMetadataProvider) IDs() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	ids := make([]string, 0, len(p.agents))
	for id := range p.agents {
		ids = append(ids, id)
	}
	return ids
}

func (p *AgentMetadataProvider) Add(ctx context.Context, agent string) error {
	metadata, err := p.storage.GetMetadata(ctx, agent)

	if err != nil {
		return fmt.Errorf("failed to obtain metadata from agent: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.agents[agent] = metadata
	return nil
}

func (p *AgentMetadataProvider) Remove(agent string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.agents, agent)
	return nil
}
