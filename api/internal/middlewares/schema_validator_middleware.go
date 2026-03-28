package middlewares

import (
	"log/slog"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"zhacked.me/oxyl/api/internal/models"
)

type SchemaValidator struct {
	typeData models.ValidationType
	body     reflect.Type
}

// https://docs.gofiber.io/blog/fiber-v3-binding-in-practice/
// Fiber does in situs content validation, so we just need to have a middleware that will validate the schema for us.
// And save on the locals the data so we can use it later on our controller.

func NewSchemaValidator(typeData models.ValidationType, body interface{}) *SchemaValidator {
	if body == nil {
		return nil
	}

	return &SchemaValidator{
		typeData: typeData,
		body:     reflect.TypeOf(body),
	}
}

func (s *SchemaValidator) Handle(ctx fiber.Ctx) error {
	data := reflect.New(s.body).Interface()

	var err error

	switch s.typeData {
	case models.URIData:
		err = ctx.Bind().URI(data)
	case models.JSONData:
		err = ctx.Bind().JSON(data)
	case models.QueryData:
		err = ctx.Bind().Query(data)
	case models.MixedData:
		err = ctx.Bind().All(data)
	default:
		return fiber.ErrUnsupportedMediaType
	}

	if err != nil {
		slog.Error("failed to parse data from request",
			slog.String("error", err.Error()),
			slog.String("type", string(s.typeData)),
			slog.String("body", s.body.String()),
			slog.String("request-id", ctx.Get("request-id")),
		)

		return fiber.ErrUnprocessableEntity
	}

	ctx.Locals(s.typeData, data)
	return ctx.Next()
}
