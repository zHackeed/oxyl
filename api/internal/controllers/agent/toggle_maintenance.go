package agent

import (
	"github.com/gofiber/fiber/v3"
	bind "github.com/idan-fishman/fiber-bind"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
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
	return "/agent/:agent_id/maintenance"
}

func (t *ToggleMaintenanceController) GetRequestModel() interface{} {
	return requests.UpdateAgentStatusRequest{}
}

func (t *ToggleMaintenanceController) Handle(ctx fiber.Ctx) error {
	agentID := ctx.Params("agent_id")
	if agentID == "" {
		return fiber.ErrBadRequest
	}

	request, ok := ctx.Locals(bind.JSON).(*requests.UpdateAgentStatusRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := t.agentService.UpdateAgentStatus(ctx.Context(), agentID, request.Status); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "maintenance toggled"})
}
