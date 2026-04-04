package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*AgentLoginController)(nil)

type AgentLoginController struct {
	agentService *service.AgentService
	tokenService *service.TokenService
}

func NewAgentLoginController(agentService *service.AgentService, tokenService *service.TokenService) *AgentLoginController {
	return &AgentLoginController{
		tokenService: tokenService,
		agentService: agentService,
	}
}

func (a AgentLoginController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (a AgentLoginController) GetPath() string {
	return "/agent/auth/login"
}

func (a AgentLoginController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.AgentLoginRequest{})
}

func (a AgentLoginController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(a.RequestRequirements().GetValidationType()).(*requests.AgentLoginRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if len(request.AgentId) < 26 {
		return fiber.ErrUnauthorized
	}

	internalCtx := context.WithValue(ctx, models.ContextInternal, true)
	agent, err := a.agentService.GetAgent(internalCtx, request.AgentId)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPermissionDenied):
			return fiber.ErrUnauthorized
		case errors.Is(err, storage.ErrAgentNotFound):
			return fiber.ErrNotFound
		default:
			slog.Error("failed to get agent", "error", err)
			return fiber.ErrInternalServerError
		}
	}

	/*
		ip := net.ParseIP(ctx.IP())

		if !agent.RegisteredIP.Equal(ip) {
			return fiber.ErrUnauthorized
		}
	*/
	tokens, err := a.tokenService.CreateToken(agent.ID, &agent.Holder, models.TokenTypeAgent)
	if err != nil {
		slog.Error("failed to create token", "error", err)
		return fiber.ErrInternalServerError
	}

	err = a.agentService.UpdateAgentStatus(internalCtx, request.AgentId, models.AgentStatusActive)
	if err != nil {
		slog.Error("failed to update agent status", "error", err)
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(tokens)
}
