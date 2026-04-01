package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*AgentRefreshController)(nil)

type AgentRefreshController struct {
	tokenService *service.TokenService
}

func NewAgentRefreshController(tokenService *service.TokenService) *AgentRefreshController {
	return new(AgentRefreshController{
		tokenService: tokenService,
	})
}

func (a AgentRefreshController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (a AgentRefreshController) GetPath() string {
	return "agent/auth/refresh"
}

func (a AgentRefreshController) RequestRequirements() *apiModel.RequestRequirements {
	return nil
}

func (a AgentRefreshController) Handle(ctx fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// authorization: Bearer XXX....
	if !strings.HasPrefix(token, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	token = token[len("Bearer "):] // Strip "Bearer "

	tokens, err := a.tokenService.RefreshToken(ctx, token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return ctx.Status(fiber.StatusOK).JSON(tokens)
}
