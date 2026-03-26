package agent

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*ListAgentsController)(nil)

type ListAgentsController struct {
	agentService *service.AgentService
}

func NewListAgentsController(agentService *service.AgentService) *ListAgentsController {
	return &ListAgentsController{
		agentService: agentService,
	}
}

func (l *ListAgentsController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (l *ListAgentsController) GetPath() string {
	return "/company/:id/agent"
}

func (l *ListAgentsController) GetRequestModel() interface{} {
	return nil
}

func (l *ListAgentsController) Handle(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	if companyID == "" {
		return fiber.ErrBadRequest
	}

	agents, err := l.agentService.GetAgents(ctx.Context(), companyID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(agents)
}
