package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	request "zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*AgentShutdownController)(nil)

type AgentShutdownController struct {
	agentService *service.AgentService
	tokenService *service.TokenService
}

func NewAgentShutdownController(agentService *service.AgentService, tokenService *service.TokenService) *AgentShutdownController {
	return new(AgentShutdownController{
		agentService: agentService,
		tokenService: tokenService,
	})
}

func (a *AgentShutdownController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (a *AgentShutdownController) GetPath() string {
	return "/agent/auth/shutdown"
}

func (a *AgentShutdownController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, request.AuthenticationShutdownRequest{})
}

func (a *AgentShutdownController) Handle(ctx fiber.Ctx) error {
	shutdownRequest, ok := ctx.Locals(a.RequestRequirements().GetValidationType()).(*request.AuthenticationShutdownRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	token := ctx.Get("Authorization")
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// authorization: Bearer XXX....
	if !strings.HasPrefix(token, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	token = token[len("Bearer "):] // Strip "Bearer "

	if err := a.tokenService.RevokeToken(ctx, token); err != nil {
		return fiber.ErrUnauthorized
	}

	internalCtx := context.WithValue(ctx, models.ContextInternal, true)
	if err := a.agentService.UpdateAgentStatus(internalCtx, shutdownRequest.AgentId, models.AgentStatusInactive); err != nil {
		slog.Error("Failed to update agent status", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(http.StatusNoContent)
}
