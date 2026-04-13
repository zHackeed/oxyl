package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var (
	ErrInvalidAuthorizationHeader = status.Error(codes.Unauthenticated, "invalid authorization header")
	ErrMissingAuthorizationHeader = status.Error(codes.Unauthenticated, "missing authorization header")
	ErrInvalidToken               = status.Error(codes.Unauthenticated, "invalid token")
)

type AgentAuthInterceptor struct {
	tokenService *service.TokenService
	agentService *service.AgentService
}

func NewAgentAuthInterceptor(tokenService *service.TokenService, agentService *service.AgentService) *AgentAuthInterceptor {
	return &AgentAuthInterceptor{
		tokenService: tokenService,
		agentService: agentService,
	}
}

func (a *AgentAuthInterceptor) Intercept(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingAuthorizationHeader
	}

	authHeader := md.Get("authorization")
	if len(authHeader) < 1 {
		return nil, ErrMissingAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader[0], "Bearer ") {
		return nil, ErrInvalidAuthorizationHeader
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == "" {
		return nil, ErrInvalidAuthorizationHeader
	}

	parsedToken, err := a.tokenService.ParseToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if parsedToken.Type != models.TokenTypeAgent {
		return nil, ErrInvalidToken
	}

	internalCtx := context.WithValue(ctx, models.ContextKeyAgent, parsedToken.Identifier)
	return handler(internalCtx, req)
}
