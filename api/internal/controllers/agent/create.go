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

var _ apiModel.Registrable = (*CreateAgentController)(nil)

type CreateAgentController struct {
	agentService *service.AgentService
}

func NewCreateAgentController(agentService *service.AgentService) *CreateAgentController {
	return &CreateAgentController{
		agentService: agentService,
	}
}

func (c *CreateAgentController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (c *CreateAgentController) GetPath() string {
	return "/agent/register"
}

func (c *CreateAgentController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.CreateAgentRequest{})
}

func (c *CreateAgentController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(c.RequestRequirements().GetValidationType()).(*requests.CreateAgentRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	agent, err := c.agentService.CreateAgent(ctx, request.DisplayName, request.RegisteredIP, request.Holder)
	if err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		slog.Error("unable to create agent", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusCreated).JSON(agent)
}
