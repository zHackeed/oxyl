package agent

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
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

func (d *DeleteAgentController) GetRequestModel() interface{} {
	return nil
}

func (d *DeleteAgentController) Handle(ctx fiber.Ctx) error {
	agentID := ctx.Params("agent_id")
	if agentID == "" {
		return fiber.ErrBadRequest
	}
	if err := d.agentService.DeleteAgent(ctx.Context(), agentID); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
