package agent

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
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
	return "/company/:company_id/agents"
}

func (l *ListAgentsController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (l *ListAgentsController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(l.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	agents, err := l.agentService.GetAgents(ctx, request.CompanyId)
	if err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}

		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(agents)
}
