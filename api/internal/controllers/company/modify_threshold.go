package company

import (
	"github.com/gofiber/fiber/v3"
	bind "github.com/idan-fishman/fiber-bind"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
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

func (m *ModifyThresholdController) GetRequestModel() interface{} {
	return requests.ModifyThresholdRequest{}
}

func (m *ModifyThresholdController) Handle(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	if companyID == "" {
		return fiber.ErrBadRequest
	}

	request, ok := ctx.Locals(bind.JSON).(*requests.ModifyThresholdRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := m.companyService.SetNotificationThreshold(ctx.Context(), companyID, request.NotificationType, request.Threshold); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "threshold modified"})
}
