package agent

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*ToggleMaintenanceController)(nil)

type ToggleMaintenanceController struct {
	agentService *service.AgentService
}

func NewToggleMaintenanceController(agentService *service.AgentService) *ToggleMaintenanceController {
	return &ToggleMaintenanceController{
		agentService: agentService,
	}
}

func (t *ToggleMaintenanceController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPatch
}

func (t *ToggleMaintenanceController) GetPath() string {
	return "/agent/:id/maintenance"
}

func (t *ToggleMaintenanceController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.MixedData, requests.UpdateAgentStatusRequest{})
}

func (t *ToggleMaintenanceController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(t.RequestRequirements().GetValidationType()).(*requests.UpdateAgentStatusRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := t.agentService.UpdateAgentStatus(ctx, request.Agent, request.Status); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		slog.Error("unable to update agent status", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "maintenance toggled"})
}
