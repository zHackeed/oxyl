package company

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

var _ apiModel.Registrable = (*AddMemberController)(nil)

type AddMemberController struct {
	companyService *service.CompanyService
}

func NewAddMemberController(companyService *service.CompanyService) *AddMemberController {
	return &AddMemberController{
		companyService: companyService,
	}
}

func (a *AddMemberController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (a *AddMemberController) GetPath() string {
	return "/company/:id/member"
}

func (a *AddMemberController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.MixedData, requests.AddMemberRequest{})
}

func (a *AddMemberController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(a.RequestRequirements().GetValidationType()).(*requests.AddMemberRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := a.companyService.AddUserToCompany(ctx, request.CompanyId, request.UserEmail, int(request.Permission)); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		if errors.Is(err, storage.ErrMemberAlreadyExists) {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "member already exists"})
		}

		slog.Error("unable to add member to company", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member added"})
}
