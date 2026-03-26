package middlewares

import (
	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*AuthMiddleware)(nil)

type AuthMiddleware struct {
	tokenService *service.TokenService
}

func NewAuthMiddleware(tokenService *service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (a *AuthMiddleware) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodNone
}

func (a *AuthMiddleware) GetPath() string {
	// the path is not relevant for a middleware, must englobe all the routes.
	return ""
}

func (a *AuthMiddleware) GetRequestModel() interface{} {
	return nil
}

func (a *AuthMiddleware) Handle(ctx fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// authorization: Bearer XXX....
	token = token[7:] // Strip "Bearer "

	parsedToken, err := a.tokenService.ParseToken(token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if parsedToken.Type != models.TokenTypeUser {
		return fiber.ErrUnauthorized
	}

	ctx.Locals(models.ContextKeyUser, parsedToken.Identifier)
	return ctx.Next()
}
