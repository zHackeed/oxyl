package company

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

type InfoController struct {
	companyService *service.CompanyService
}

var _ apiModel.Registrable = (*InfoController)(nil)

func NewInfoController(companyService *service.CompanyService) *InfoController {
	return &InfoController{
		companyService: companyService,
	}
}

func (i *InfoController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (i *InfoController) GetPath() string {
	return "/company/:id"
}

func (i *InfoController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (i *InfoController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(i.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	company, err := i.companyService.GetCompany(ctx, request.CompanyId)
	if err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(company)
}
