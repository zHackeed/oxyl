package company

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*ThresholdsController)(nil)

type ThresholdsController struct {
	companyService *service.CompanyService
}

func NewThresholdsController(companyService *service.CompanyService) *ThresholdsController {
	return &ThresholdsController{
		companyService: companyService,
	}
}

func (t *ThresholdsController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (t *ThresholdsController) GetPath() string {
	return "/company/:id/thresholds"
}

func (t *ThresholdsController) GetRequestModel() interface{} {
	return nil
}

func (t *ThresholdsController) Handle(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	if companyID == "" {
		return fiber.ErrBadRequest
	}

	thresholds, err := t.companyService.GetNotificationThresholds(ctx.Context(), companyID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(thresholds)
}
