package member

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*RemoveMemberController)(nil)

type RemoveMemberController struct {
	companyService *service.CompanyService
	userService    *service.UserService
}

func NewRemoveMemberController(companyService *service.CompanyService, userService *service.UserService) *RemoveMemberController {
	return &RemoveMemberController{
		companyService: companyService,
		userService:    userService,
	}
}

func (r *RemoveMemberController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodDelete
}

func (r *RemoveMemberController) GetPath() string {
	return "/company/:company_id/member/"
}

func (r *RemoveMemberController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.MixedData, requests.RemoveMemberRequest{})
}

func (r *RemoveMemberController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(r.RequestRequirements().GetValidationType()).(*requests.RemoveMemberRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	userId, err := r.userService.GetIdFromEmail(ctx, request.UserEmail)
	if err != nil {
		return fiber.ErrNotFound
	}

	if err := r.companyService.RemoveUserFromCompany(ctx, request.CompanyId, userId); err != nil {
		switch {
		case errors.Is(err, storage.ErrMemberNotFound):
			return fiber.ErrNotFound
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		default:
			slog.Error("unable to remove member from company", "error", err)
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
