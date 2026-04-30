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

var _ apiModel.Registrable = (*CreateEntrypointController)(nil)

type CreateEntrypointController struct {
	companyService *service.CompanyService
}

func NewCreateEntrypointController(agentService *service.CompanyService) *CreateEntrypointController {
	return &CreateEntrypointController{companyService: agentService}
}

func (c *CreateEntrypointController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (c *CreateEntrypointController) GetPath() string {
	return "/company/:id/notification/entrypoint"
}

func (c *CreateEntrypointController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.MixedData, requests.CreateEndpointRequest{})
}

func (c *CreateEntrypointController) Handle(ctx fiber.Ctx) error {
	request, allowed := ctx.Locals(c.RequestRequirements().GetValidationType()).(*requests.CreateEndpointRequest)
	if !allowed {
		slog.Info("request not allowed")
		return fiber.ErrUnauthorized
	}

	endpoint := &models.CompanyNotificationSettings{
		WebhookType: request.WebhookType,
		Endpoint:    request.Endpoint,
	}

	if request.Channel != nil {
		endpoint.Channel = request.Channel
	}

	notificationSetting, err := c.companyService.AddNotificationEndpoint(ctx, request.CompanyId, endpoint)

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
