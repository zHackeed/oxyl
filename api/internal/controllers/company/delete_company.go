package company

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*DeleteCompanyController)(nil)

type DeleteCompanyController struct {
	companyService *service.CompanyService
}

func NewDeleteCompanyController(companyService *service.CompanyService) *DeleteCompanyController {
	return &DeleteCompanyController{
		companyService: companyService,
	}
}

func (d *DeleteCompanyController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodDelete
}

func (d *DeleteCompanyController) GetPath() string {
	return "/company/:id"
}

func (d *DeleteCompanyController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (d *DeleteCompanyController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(d.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := d.companyService.Delete(ctx, request.CompanyId); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		slog.Error("unable to delete company", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "company deleted"})
}
