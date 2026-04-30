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

var _ apiModel.Registrable = (*DeleteEntrypointController)(nil)

type DeleteEntrypointController struct {
	companyService *service.CompanyService
}

func NewDeleteEntrypointController(agentService *service.CompanyService) *CreateEntrypointController {
	return &CreateEntrypointController{companyService: agentService}
}

func (d *DeleteEntrypointController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodDelete
}

func (d *DeleteEntrypointController) GetPath() string {
	return "/company/:id/notification/entrypoint/:entrypoint"
}

func (d *DeleteEntrypointController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.DeleteEndpointRequest{})
}

func (d *DeleteEntrypointController) Handle(ctx fiber.Ctx) error {
	request, allowed := ctx.Locals(d.RequestRequirements().GetValidationType()).(*requests.DeleteEndpointRequest)
	if !allowed {
		return fiber.ErrUnauthorized
	}

	err := d.companyService.RemoveNotificationEndpoint(ctx, request.CompanyId, request.EndpointId)

	if err != nil {
		slog.Info(err.Error())
		switch {
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		case errors.Is(err, storage.ErrNotificationEndpointNotFound):
			return fiber.ErrNotFound
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
