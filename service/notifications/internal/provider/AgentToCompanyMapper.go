package provider

import (
	"context"
	"sync"

	"zhacked.me/oxyl/service/notifications/internal/storage"
)

type AgentCompanyProvider struct {
	mu     sync.RWMutex
	agents map[string]string // agentID -> companyID

	storage *storage.AgentToCompanyMapperStorage
}

func NewAgentCompanyProvider(storage *storage.AgentToCompanyMapperStorage) *AgentCompanyProvider {
	return &AgentCompanyProvider{
		agents:  make(map[string]string),
		storage: storage,
	}
}

func (p *AgentCompanyProvider) Load(ctx context.Context) error {
	agents, err := p.storage.GetAll(ctx)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.agents = agents
	return nil
}

func (p *AgentCompanyProvider) Get(agentID string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	companyID, ok := p.agents[agentID]
	return companyID, ok
}

func (p *AgentCompanyProvider) Add(ctx context.Context, agentID string) error {
	companyId, err := p.storage.GetByAgent(ctx, agentID)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.agents[agentID] = companyId
	return nil
}

func (p *AgentCompanyProvider) Remove(agentID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.agents, agentID)
}
