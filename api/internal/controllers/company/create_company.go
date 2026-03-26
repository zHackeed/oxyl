package company

import (
	"github.com/gofiber/fiber/v3"
	bind "github.com/idan-fishman/fiber-bind"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*CreateCompanyController)(nil)

type CreateCompanyController struct {
	companyService *service.CompanyService
}

func NewCreateCompanyController(companyService *service.CompanyService) *CreateCompanyController {
	return &CreateCompanyController{
		companyService: companyService,
	}
}

func (c *CreateCompanyController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (c *CreateCompanyController) GetPath() string {
	return "/company/create"
}

func (c *CreateCompanyController) GetRequestModel() interface{} {
	return requests.CreateCompanyRequest{}
}

func (c *CreateCompanyController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(bind.JSON).(*requests.CreateCompanyRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	_, err := c.companyService.CreateCompany(ctx.Context(), request.DisplayName)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "company created"})
}
