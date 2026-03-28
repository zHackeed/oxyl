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

var _ apiModel.Registrable = (*DeleteAgentController)(nil)

type DeleteAgentController struct {
	agentService *service.AgentService
}

func NewDeleteAgentController(agentService *service.AgentService) *DeleteAgentController {
	return &DeleteAgentController{
		agentService: agentService,
	}
}

func (d *DeleteAgentController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodDelete
}

func (d *DeleteAgentController) GetPath() string {
	return "/agent/:agent_id"
}

func (d *DeleteAgentController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.AgentIdUri{})
}

func (d *DeleteAgentController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(d.RequestRequirements().GetValidationType()).(*requests.AgentIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := d.agentService.DeleteAgent(ctx, request.AgentId); err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		slog.Error("unable to delete agent", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
