package user

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
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
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.RefreshTokenRequest{})
}

func NewRefreshController(tokenService *service.TokenService) *RefreshController {
	return &RefreshController{
		tokenService: tokenService,
	}
}

func (r *RefreshController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(r.RequestRequirements().GetValidationType()).(*requests.RefreshTokenRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	tokenPair, err := r.tokenService.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		slog.Error("unable to refresh token", "error", err)
		return fiber.ErrUnauthorized
	}

	return ctx.Status(fiber.StatusOK).JSON(tokenPair)
}
