package user

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
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
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.LogoutRequest{})
}

func (l LogoutController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(l.RequestRequirements().GetValidationType()).(*requests.LogoutRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := l.tokenService.RevokeToken(ctx, request.RefreshToken); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout successful"})
}
