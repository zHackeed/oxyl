package user

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*RefreshController)(nil)

type RefreshController struct {
	tokenService *service.TokenService
}

func (r *RefreshController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (r *RefreshController) GetPath() string {
	return "/auth/refresh"
}

func (r *RefreshController) RequestRequirements() *apiModel.RequestRequirements {
	return nil
}

func NewRefreshController(tokenService *service.TokenService) *RefreshController {
	return &RefreshController{
		tokenService: tokenService,
	}
}

func (r *RefreshController) Handle(ctx fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// authorization: Bearer XXX....
	if !strings.HasPrefix(token, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	token = token[len("Bearer "):] // Strip "Bearer "

	tokenPair, err := r.tokenService.RefreshToken(ctx, token)
	if err != nil {
		slog.Error("unable to refresh token", "error", err)
		return fiber.ErrUnauthorized
	}

	return ctx.Status(fiber.StatusOK).JSON(tokenPair)
}
