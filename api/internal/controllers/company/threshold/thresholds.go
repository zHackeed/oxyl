package threshold

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/api/internal/models/responses"
	"zhacked.me/oxyl/shared/pkg/models"
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

func (t *ThresholdsController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (t *ThresholdsController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(t.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	thresholds, err := t.companyService.GetNotificationThresholds(ctx, request.CompanyId)
	if err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		slog.Error("unable to get company thresholds", "error", err)
		return fiber.ErrInternalServerError
	}

	valueWrapped := make([]responses.CompanyThresholdValueWrapper, len(thresholds))

	for key, value := range thresholds {
		valueWrapped = append(valueWrapped, responses.CompanyThresholdValueWrapper{
			ThresholdIdentifier: key,
			Value:               value,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(valueWrapped)
}
