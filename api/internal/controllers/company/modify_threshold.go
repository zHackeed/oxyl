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

var _ apiModel.Registrable = (*ModifyThresholdController)(nil)

type ModifyThresholdController struct {
	companyService *service.CompanyService
}

func NewModifyThresholdController(companyService *service.CompanyService) *ModifyThresholdController {
	return &ModifyThresholdController{
		companyService: companyService,
	}
}

func (m *ModifyThresholdController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPatch
}

func (m *ModifyThresholdController) GetPath() string {
	return "/company/:id/thresholds"
}

func (m *ModifyThresholdController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.ModifyThresholdRequest{})
}

func (m *ModifyThresholdController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(m.RequestRequirements().GetValidationType()).(*requests.ModifyThresholdRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := m.companyService.SetNotificationThreshold(ctx, request.CompanyId, request.NotificationType, request.Threshold); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		slog.Error("unable to modify threshold", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "threshold modified"})
}
