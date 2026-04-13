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

// TODO: might think of caching company user lookups?

type CompanyService struct {
	messenger      *datasource.RedisConnection
	userStorage    *storage.UserStorage
	companyStorage *storage.CompanyStorage
}

func NewCompanyService(
	messenger *datasource.RedisConnection,
	companyStorage *storage.CompanyStorage,
	userStorage *storage.UserStorage,
) *CompanyService {
	return &CompanyService{
		messenger:      messenger,
		companyStorage: companyStorage,
		userStorage:    userStorage,
	}
}

func (c *CompanyService) CreateCompany(ctx context.Context, displayName string) (*models.Company, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	model, err := models.NewCompany(displayName, userId)
	if err != nil {
		return nil, fmt.Errorf("unable to create company: %w", err)
	}

	if err := c.companyStorage.CreateCompany(ctx, model); err != nil {
		return nil, fmt.Errorf("unable to create company: %w", err)
	}

	return model, nil
}

func (c *CompanyService) GetCompany(ctx context.Context, companyId string) (*models.Company, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionView)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	return c.companyStorage.GetCompany(ctx, companyId)
}

func (c *CompanyService) GetCompanies(ctx context.Context) ([]*models.Company, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	return c.companyStorage.GetCompaniesOfUser(ctx, userId)
}

func (c *CompanyService) GetMember(ctx context.Context, companyId string) (*models.CompanyMember, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	member, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (c *CompanyService) AddUserToCompany(ctx context.Context, companyId, userEmail string, permission int) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return models.ErrPermissionDenied
	}
	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageMembers)
	if !allowed {
		return models.ErrPermissionDenied
	}

	targetId, err := c.userStorage.GetIdFromEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("unable to get user: %w", err)
	}

	err = c.companyStorage.AddMemberToCompany(ctx, targetId, companyId, permission)
	if err != nil {
		if errors.Is(err, storage.ErrMemberAlreadyExists) {
			// propagate the same error
			return storage.ErrMemberAlreadyExists
		}

		return fmt.Errorf("unable to add member to company: %w", err)
	}

	if err := c.messenger.Publish(ctx, variables.RedisChannelCompanyAddedMember, redisModels.CompanyMemberAdded{
		CompanyId: companyId,
		UserId:    targetId,
	}); err != nil {
		return fmt.Errorf("unable to publish company update event: %w", err)
	}

	return nil
}

func (c *CompanyService) GetMembers(ctx context.Context, companyId string) ([]*models.CompanyMemberComposite, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}
	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			return nil, models.ErrPermissionDenied
		}

		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageMembers)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	members, err := c.companyStorage.GetMembersOfCompany(ctx, companyId)
	if err != nil {
		return nil, err // explode
	}

	return members, nil
}

func (c *CompanyService) RemoveUserFromCompany(ctx context.Context, companyId, targetId string) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return models.ErrPermissionDenied
	}

	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageMembers)
	if !allowed {
		return models.ErrPermissionDenied
	}

	targetMembership, err := c.companyStorage.GetCompanyMembership(ctx, targetId, companyId)
	if err != nil {
		return err
	}

	if targetMembership.Permission == models.CompanyPermissionOwner {
		return models.ErrPermissionDenied
	}

	if err := c.companyStorage.RemoveMemberFromCompany(ctx, companyId, targetId); err != nil {
		return fmt.Errorf("unable to remove member from company: %w", err)
	}

	if err := c.messenger.Publish(ctx, variables.RedisChannelCompanyRemovedMember, redisModels.CompanyMemberRemoved{
		CompanyId: companyId,
		UserId:    targetId,
	}); err != nil {
		slog.Error("unable to publish company update event", "error", err)
	}

	return nil
}

func (c *CompanyService) GetNotificationThresholds(ctx context.Context, companyId string) (map[models.NotificationType]int, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, models.ErrPermissionDenied
	}

	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return nil, err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageThresholds)
	if !allowed {
		return nil, models.ErrPermissionDenied
	}

	return c.companyStorage.GetNotificationThresholdsOfCompany(ctx, companyId)
}

func (c *CompanyService) SetNotificationThreshold(ctx context.Context, companyId string, notificationType models.NotificationType, threshold int) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return models.ErrPermissionDenied
	}

	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionManageThresholds)
	if !allowed {
		return models.ErrPermissionDenied
	}

	if err := c.companyStorage.SetNotificationThresholdForCompany(ctx, companyId, notificationType, threshold); err != nil {
		return fmt.Errorf("unable to set notification threshold for company: %w", err)
	}

	// cache invalidation
	if err := c.messenger.Publish(ctx, variables.RedisChannelCompanyThresholdUpdate, redisModels.ThresholdUpdate{
		CompanyId:     companyId,
		ThresholdType: notificationType,
		Threshold:     threshold,
	}); err != nil {
		slog.Error("unable to publish company update event", "error", err)
	}

	return nil
}

func (c *CompanyService) Delete(ctx context.Context, companyId string) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return models.ErrPermissionDenied
	}

	if companyId == "" {
		return errors.New("company id is empty")
	}

	if len(companyId) > 26 {
		return errors.New("company id is too long, maybe malformed")
	}

	membership, err := c.companyStorage.GetCompanyMembership(ctx, userId, companyId)
	if err != nil {
		return err
	}

	allowed := models.HasPermission(membership.Permission, models.CompanyPermissionOwner) // Only owners can delete companies.
	if !allowed {
		return models.ErrPermissionDenied
	}

	if err := c.companyStorage.DeleteCompany(ctx, companyId); err != nil {
		return fmt.Errorf("unable to delete company: %w", err)
	}

	if err := c.messenger.Publish(ctx, variables.RedisChannelCompanyDeletion, redisModels.CompanyDeletion{
		CompanyId: companyId,
	}); err != nil {
		slog.Error("unable to publish company update event", "error", err)
	}

	return nil
}

/*
	Todo: notification endpoint management
*/
