package notifications

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*ListController)(nil)

type ListController struct {
	agentService *service.AgentService
}

func NewListController(agentService *service.AgentService) *ListController {
	return &ListController{
		agentService: agentService,
	}
}

func (l ListController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (l ListController) GetPath() string {
	return "/agent/:id/notifications"
}

func (l ListController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.AgentIdUri{})
}

func (l ListController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(l.RequestRequirements().GetValidationType()).(*requests.AgentIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	logs, err := l.agentService.GetNotificationLogs(ctx, request.AgentId)

	if err != nil {
		switch {
		case errors.Is(err, storage.ErrAgentNotEnrolled):
			return fiber.ErrExpectationFailed
		case errors.Is(err, storage.ErrAgentNotFound):
			return fiber.ErrNotFound
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		}
		return err
	}

	return ctx.JSON(logs)

}
