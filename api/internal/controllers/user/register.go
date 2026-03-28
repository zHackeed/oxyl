package user

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*RegisterController)(nil)

type RegisterController struct {
	userService *service.UserService
}

func NewRegisterController(userService *service.UserService) *RegisterController {
	return &RegisterController{
		userService: userService,
	}
}

func (r *RegisterController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (r *RegisterController) GetPath() string {
	return "/auth/register"
}

func (r *RegisterController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.JSONData, requests.RegisterRequest{})
}

func (r *RegisterController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(r.RequestRequirements().GetValidationType()).(*requests.RegisterRequest)
	if !ok {
		return fiber.ErrInternalServerError
	}

	_, err := r.userService.Register(ctx, request.Name, request.Surname, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return fiber.ErrConflict // 409 - a user with that email already exists.
		}

		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
}
