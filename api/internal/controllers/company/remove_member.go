package company

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*RemoveMemberController)(nil)

type RemoveMemberController struct {
	companyService *service.CompanyService
}

func NewRemoveMemberController(companyService *service.CompanyService) *RemoveMemberController {
	return &RemoveMemberController{
		companyService: companyService,
	}
}

func (r *RemoveMemberController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodDelete
}

func (r *RemoveMemberController) GetPath() string {
	return "/company/:company_id/member/:user_id"
}

func (r *RemoveMemberController) GetRequestModel() interface{} {
	return nil
}

func (r *RemoveMemberController) Handle(ctx fiber.Ctx) error {
	companyID := ctx.Params("company_id")
	if companyID == "" {
		return fiber.ErrBadRequest
	}

	userID := ctx.Params("user_id")
	if userID == "" {
		return fiber.ErrBadRequest
	}

	if err := r.companyService.RemoveUserFromCompany(ctx.Context(), companyID, userID); err != nil {
		switch {
		case errors.Is(err, storage.ErrMemberNotFound):
			return fiber.ErrNotFound
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
