package agent

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
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
		switch {
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		case errors.Is(err, storage.ErrNoAgents):
			return fiber.ErrNotFound
		default:
			slog.Error("failed to get agents for company", slog.Any("error", err), slog.String("company_id", request.CompanyId))
			return fiber.ErrInternalServerError
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(agents)
}
