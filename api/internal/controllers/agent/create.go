package agent

import (
	"github.com/gofiber/fiber/v3"
	bind "github.com/idan-fishman/fiber-bind"
	"zhacked.me/oxyl/api/internal/models"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
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

func (c *CreateAgentController) GetMethod() models.HttpMethod {
	return models.MethodPost
}

func (c *CreateAgentController) GetPath() string {
	return "/agent/register"
}

func (c *CreateAgentController) GetRequestModel() interface{} {
	return requests.CreateAgentRequest{}
}

func (c *CreateAgentController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(bind.JSON).(*requests.CreateAgentRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	_, err := c.agentService.CreateAgent(ctx.Context(), request.DisplayName, request.RegisteredIP, request.Holder)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "agent created"})
}
