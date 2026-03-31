package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"zhacked.me/oxyl/shared/pkg/datasource"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/utils"
	"zhacked.me/oxyl/shared/pkg/variables"
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

func (a *AgentService) CreateAgent(ctx context.Context, displayName, registeredIP, holder string) (*models.Agent, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, holder)
	if err != nil {
		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageAgents)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	model, err := models.NewAgent(displayName, registeredIP, holder)
	if err != nil {
		return nil, fmt.Errorf("unable to create agent: %w", err)
	}

	//todo: agent company limit count handling

	if err := a.agentStorage.CreateAgent(ctx, model); err != nil {
		return nil, fmt.Errorf("unable to save agent to storage: %w", err)
	}

	// notify the users watching the interface and resort the gui on their end.
	if err := a.messenger.Publish(ctx, variables.RedisChannelAgentUpdate, redisModels.AgentCreation{
		CompanyId:    holder,
		AgentId:      model.ID,
		State:        models.AgentStatusEnrolling,
		RegisteredIP: registeredIP,
		DisplayName:  displayName,
	}); err != nil {
		slog.Error("unable to publish agent creation event", "error", err)
	}

	return model, nil
}

func (a *AgentService) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	userId, userFound := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	agentId, agentFound := utils.GetValueFromContext[string](ctx, models.ContextAgent)
	_, internalCall := utils.GetValueFromContext[bool](ctx, models.ContextInternal)

	if !userFound && !agentFound && !internalCall {
		return nil, models.ErrPermissionDenied
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

	if userFound {
		membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, agent.Holder)
		if err != nil {
			return nil, err
		}

		if !models.HasPermission(membership.Permission, models.CompanyPermissionView) {
			return nil, models.ErrPermissionDenied
		}
	}

	if agentFound && agentId != agentID {
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
	userId, userFound := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	_, internalCall := utils.GetValueFromContext[bool](ctx, models.ContextInternal)

	if !userFound && !internalCall {
		return models.ErrPermissionDenied
	}

	agent, err := a.agentStorage.GetAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("unable to get agent: %w", err)
	}

	if userFound {
		membership, err := a.companyStorage.GetCompanyMembership(ctx, userId, agent.Holder)
		if err != nil {
			return fmt.Errorf("unable to get company membership: %w", err)
		}

		if !models.HasPermission(membership.Permission, models.CompanyPermissionManageAgents) {
			return models.ErrPermissionDenied
		}
	}

	// notify agent of status update
	if err := a.messenger.Publish(ctx, variables.RedisChannelAgentUpdate, redisModels.AgentUpdate{
		AgentId: agentID,
		Status:  status,
	}); err != nil {
		slog.Error("unable to publish agent update event", "error", err)
	}

	return nil
}

func (a *AgentService) EnrichAgent(ctx context.Context, agentID string,
	osName, cpuModel string,
	totalMemory, totalDisk int64,
	partitions []*models.AgentPartition,
	enrollmentToken string,
) error {
	_, internalCall := utils.GetValueFromContext[bool](ctx, models.ContextInternal)
	if !internalCall {
		return models.ErrPermissionDenied
	}

	if err := a.agentStorage.EnrichAgent(ctx,
		agentID, osName, cpuModel,
		totalMemory, totalDisk,
		enrollmentToken, partitions,
	); err != nil {
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

	// We need to disconnect the agent (kill it) from the gRPC and also tell the websockets to stop reading data and redirect them back to the agent list
	if err := a.messenger.Publish(ctx, variables.RedisChannelAgentDeletion, redisModels.AgentDelete{
		CompanyId: agent.Holder,
		AgentId:   agentID,
	}); err != nil {
		slog.Error("unable to publish agent deletion event", "error", err)
	}

	return nil
}
