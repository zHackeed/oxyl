package service

import (
	"context"
	"errors"
	"fmt"

	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/utils"
)

type AgentService struct {
	messenger *datasource.RedisConnection

	companyStorage *storage.CompanyStorage
	agentStorage   *storage.AgentStorage
}

func NewAgentService(
	messenger *datasource.RedisConnection,
	companyStorage *storage.CompanyStorage,
	agentStorage *storage.AgentStorage,
) *AgentService {
	return &AgentService{
		messenger:      messenger,
		companyStorage: companyStorage,
		agentStorage:   agentStorage,
	}
}

func (a *AgentService) CreateAgent(ctx context.Context, displayNam, registeredIP, holder string) (*models.Agent, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, errors.New("user not found in context")
	}

	membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, holder)
	if err != nil {
		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageAgents)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	model, err := models.NewAgent(displayNam, registeredIP, holder)
	if err != nil {
		return nil, fmt.Errorf("unable to create agent: %w", err)
	}

	//todo: agent company limit count handling

	if err := a.agentStorage.CreateAgent(ctx, model); err != nil {
		return nil, fmt.Errorf("unable to save agent to storage: %w", err)
	}

	return model, nil
}

func (a *AgentService) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, errors.New("user not found in context")
	}

	if agentID == "" {
		return nil, errors.New("agent id is empty")
	}

	if len(agentID) > 26 {
		return nil, errors.New("agent id is too long, maybe malformed")
	}

	agent, err := a.agentStorage.GetAgent(ctx, agentID)
	if err != nil {
		return nil, err
	}

	membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, agent.Holder)
	if err != nil {
		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionView)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	return agent, nil
}

func (a *AgentService) GetAgents(ctx context.Context, companyID string) ([]*models.Agent, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, errors.New("user not found in context")
	}

	membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, companyID)
	if err != nil {
		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionView)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	return a.agentStorage.GetAgentsOfCompany(ctx, companyID)
}

func (a *AgentService) UpdateAgentStatus(ctx context.Context, agentID string, status models.AgentStatus) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return errors.New("user not found in context")
	}

	agent, err := a.agentStorage.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("unable to get agent: %w", err)
	}

	membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, agent.Holder)
	if err != nil {
		return fmt.Errorf("unable to get company membership: %w", err)
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageAgents)
	if !allowed {
		return models.ErrPermissionDenied
	}

	if err := a.agentStorage.UpdateAgentStatus(ctx, agentID, status); err != nil {
		return fmt.Errorf("unable to update agent status: %w", err)
	}

	return nil
}

func (a *AgentService) DeleteAgent(ctx context.Context, agentID string) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return errors.New("user not found in context")
	}

	agent, err := a.agentStorage.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("unable to get agent: %w", err)
	}

	membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, agent.Holder)
	if err != nil {
		return fmt.Errorf("unable to get company membership: %w", err)
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageAgents)

	if !allowed {
		return models.ErrPermissionDenied
	}

	if err := a.agentStorage.DeleteAgent(ctx, agentID); err != nil {
		return fmt.Errorf("unable to delete agent: %w", err)
	}

	return nil
}
