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

var _ apiModel.Registrable = (*InfoController)(nil)

type InfoController struct {
	agentService *service.AgentService
}

func NewAgentInfoController(agentService *service.AgentService) *InfoController {
	return &InfoController{
		agentService: agentService,
	}
}

func (a *InfoController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (a *InfoController) GetPath() string {
	return "/agent/:id"
}

func (a *InfoController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.AgentIdUri{})
}

func (a *InfoController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(a.RequestRequirements().GetValidationType()).(*requests.AgentIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	agent, err := a.agentService.GetAgent(ctx, request.AgentId)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		case errors.Is(err, storage.ErrAgentNotFound):
			return fiber.ErrNotFound
		default:
			slog.Error("unable to get agent", "error", err)
			return fiber.ErrInternalServerError
		}
	}

	return ctx.JSON(agent)
}
