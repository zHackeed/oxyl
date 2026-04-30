package endpoints

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

var _ apiModel.Registrable = (*ListEntrypointController)(nil)

type ListEntrypointController struct {
	companyService *service.CompanyService
}

func NewListEntrypointController(agentService *service.CompanyService) *CreateEntrypointController {
	return &CreateEntrypointController{companyService: agentService}
}

func (c *ListEntrypointController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (c *ListEntrypointController) GetPath() string {
	return "/company/:id/notification/entrypoint"
}

func (c *ListEntrypointController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (c *ListEntrypointController) Handle(ctx fiber.Ctx) error {
	request, allowed := ctx.Locals(c.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !allowed {
		return fiber.ErrUnauthorized
	}

	notificationSetting, err := c.companyService.GetNotificationEndpoints(ctx, request.CompanyId)

	if err != nil {
		slog.Info(err.Error())
		switch {
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		case errors.Is(err, storage.ErrCompanyNotFound):
			return fiber.ErrNotFound
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.JSON(notificationSetting)
}
