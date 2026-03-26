package company

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*ListCompaniesController)(nil)

type ListCompaniesController struct {
	companyService *service.CompanyService
}

func NewListCompaniesController(companyService *service.CompanyService) *ListCompaniesController {
	return &ListCompaniesController{
		companyService: companyService,
	}
}

func (l *ListCompaniesController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (l *ListCompaniesController) GetPath() string {
	return "/company"
}

func (l *ListCompaniesController) GetRequestModel() interface{} {
	// Get does not have a body.
	return nil
}

func (l *ListCompaniesController) Handle(ctx fiber.Ctx) error {
	companies, err := l.companyService.GetCompanies(ctx.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(companies)
}
