package agent

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*MetricsController)(nil)

type MetricsController struct {
	metricsService *service.MetricsService
}

func NewMetricsController(metricsService *service.MetricsService) *MetricsController {
	return &MetricsController{
		metricsService: metricsService,
	}
}

func (m MetricsController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (m MetricsController) GetPath() string {
	return "/agent/:id/metrics/:interval"
}

func (m MetricsController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.AgentMetricsRequest{})
}

func (m MetricsController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(m.RequestRequirements().GetValidationType()).(*requests.AgentMetricsRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	intervalParsed, err := time.ParseDuration(request.Interval)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	metrics, err := m.metricsService.GetMetrics(ctx, request.AgentId, intervalParsed)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrForbidden
		case errors.Is(err, storage.ErrAgentNotFound):
			return fiber.ErrNotFound
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.JSON(metrics)
}
