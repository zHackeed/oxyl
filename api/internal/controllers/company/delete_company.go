package company

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
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

func (d *DeleteCompanyController) GetRequestModel() interface{} {
	return nil
}

func (d *DeleteCompanyController) Handle(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	if companyID == "" {
		return fiber.ErrBadRequest
	}

	if err := d.companyService.Delete(ctx.Context(), companyID); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "company deleted"})
}
