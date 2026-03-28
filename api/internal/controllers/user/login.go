package user

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*LoginController)(nil)

type LoginController struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewLoginController(userService *service.UserService, tokenService *service.TokenService) *LoginController {
	return &LoginController{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (l *LoginController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (l *LoginController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.LoginRequest{})
}

func (l *LoginController) GetPath() string {
	return "/auth/login"
}

func (l *LoginController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(l.RequestRequirements().GetValidationType()).(*requests.LoginRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	found, err := l.userService.Authenticate(ctx, request.Email, request.Password)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	tokenPair, err := l.tokenService.CreateToken(found.ID, nil, models.TokenTypeUser)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(tokenPair)
}
