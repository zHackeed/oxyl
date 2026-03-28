package user

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*LogoutController)(nil)

type LogoutController struct {
	tokenService *service.TokenService
}

func NewLogoutController(tokenService *service.TokenService) *LogoutController {
	return &LogoutController{
		tokenService: tokenService,
	}
}

func (l LogoutController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (l LogoutController) GetPath() string {
	return "/auth/logout"
}

func (l LogoutController) RequestRequirements() *apiModel.RequestRequirements {
	return nil
}

func (l LogoutController) Handle(ctx fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// authorization: Bearer XXX....
	if !strings.HasPrefix(token, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	token = token[len("Bearer "):] // Strip "Bearer "

	if err := l.tokenService.RevokeToken(ctx, token); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout successful"})
}
