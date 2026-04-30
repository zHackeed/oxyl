package company

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
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

func (c *CreateCompanyController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.CreateCompanyRequest{})
}

func (c *CreateCompanyController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(c.RequestRequirements().GetValidationType()).(*requests.CreateCompanyRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	company, err := c.companyService.CreateCompany(ctx, request.DisplayName, request.WebhookType, request.WebhookEndpoint, request.WebhookChannel)
	if err != nil {
		slog.Error("unable to create company", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.
		Status(fiber.StatusCreated).
		JSON(company)
}
