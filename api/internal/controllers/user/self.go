package user

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*InfoController)(nil)

type InfoController struct {
	userService *service.UserService
}

func NewInfoController(userService *service.UserService) *InfoController {
	return &InfoController{
		userService: userService,
	}
}

func (i *InfoController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (i *InfoController) GetPath() string {
	return "/user/"
}

func (i *InfoController) RequestRequirements() *apiModel.RequestRequirements {
	return nil
}

func (i *InfoController) Handle(ctx fiber.Ctx) error {
	user, err := i.userService.GetUser(ctx)

	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return fiber.ErrNotFound
		default:
			slog.Error("[UserService] get user error", "error", err)
			return fiber.ErrInternalServerError
		}
	}

	return ctx.JSON(user)
}
