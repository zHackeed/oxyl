package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/utils"
)

type LoggingInterceptor struct{}

func NewLoggingInterceptor() *LoggingInterceptor {
	return &LoggingInterceptor{}
}

func (l *LoggingInterceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		slog.Error("unable to get peer from context")
	}

	agent, found := utils.GetValueFromContext[string](ctx, models.ContextKeyAgent)

	args := []any{
		slog.String("method", info.FullMethod),
		slog.String("peer", p.Addr.String()),
		//slog.Any("data", req),
	}

	if found {
		args = append(args, slog.String("agent", agent))
	}

	slog.Info("request invoked", args...)
	return handler(ctx, req)
}
