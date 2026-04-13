package middlewares

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/utils"
)

var _ apiModel.Registrable = (*LoggingMiddleware)(nil)

type LoggingMiddleware struct{}

func NewLoggingMiddleWare() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (l LoggingMiddleware) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodNone
}

func (l LoggingMiddleware) GetPath() string {
	// the path is not relevant for a middleware, must englobe all the routes.
	return ""
}

func (l LoggingMiddleware) RequestRequirements() *apiModel.RequestRequirements {
	return nil
}

func (l LoggingMiddleware) Handle(ctx fiber.Ctx) error {
	requestPath := ctx.Path()

	if requestPath == healthcheck.LivenessEndpoint {
		// Ignore
		return ctx.Next()
	}

	start := time.Now()

	userId, foundUser := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	agentId, foundAgent := utils.GetValueFromContext[string](ctx, models.ContextKeyAgent)

	method := ctx.Method()
	ip := ctx.IP()

	requestID := ctx.GetRespHeader("X-Request-ID")
	if requestID == "" {
		requestID = "none"
	}

	err := ctx.Next()
	latency := time.Since(start)
	status := ctx.Response().StatusCode()

	fields := []slog.Attr{
		slog.String("method", method),
		slog.String("path", requestPath),
		slog.String("ip", ip),
		slog.String("request_id", requestID),
		slog.Int("status", status),
		slog.String("latency", latency.String()),
	}

	if foundUser {
		fields = append(fields, slog.String("user", userId))
	}

	if foundAgent {
		fields = append(fields, slog.String("agent", agentId))
	}

	if err != nil {
		slog.Error("request failed", slog.Any("error", err), fields)
	} else {
		slog.Info("request handled", fields)
	}

	return err
}
