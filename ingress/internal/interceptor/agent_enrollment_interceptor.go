package interceptor

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/utils"
)

var (
	ErrAgentMissingHeader = status.Error(codes.InvalidArgument, "agent missing agent header")
	ErrAgentInvalidHeader = status.Error(codes.Unauthenticated, "agent invalid header")
	ErrAgentInvalidToken  = status.Error(codes.Unauthenticated, "agent invalid token")

	excludedRoutes = []string{
		"/comms.EnrollmentService/GetEnrollmentToken",
	}
)

type AgentEnrollmentInterceptor struct {
	tokenService *service.TokenService
	agentService *service.AgentService

	cache *ristretto.Cache[string, string]
}

func NewAgentEnrollmentInterceptor(tokenService *service.TokenService, agentService *service.AgentService) (*AgentEnrollmentInterceptor, error) {
	cache, err := ristretto.NewCache[string, string](&ristretto.Config[string, string]{
		NumCounters: 10000,
		MaxCost:     1000,
		BufferItems: 64,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to create cache: %w", err)
	}

	return &AgentEnrollmentInterceptor{
		tokenService: tokenService,
		agentService: agentService,
		cache:        cache,
	}, nil
}

func (a *AgentEnrollmentInterceptor) Intercept(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if slices.Contains(excludedRoutes, info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrAgentMissingHeader
	}

	agent, found := utils.GetValueFromContext[string](ctx, models.ContextAgent)
	if !found {
		return nil, ErrAgentMissingHeader
	}

	enrollmentHeader := md.Get("ag_enrollment")

	if len(enrollmentHeader) < 1 {
		return nil, ErrAgentMissingHeader
	}

	token := enrollmentHeader[0]
	if token == "" {
		return nil, ErrAgentInvalidHeader
	}

	requestedToken, err := a.requestEnrolTokenAgent(ctx, agent)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get enrollment token")
	}

	if token != requestedToken {
		return nil, ErrAgentInvalidToken
	}

	return handler(ctx, req)
}

func (a *AgentEnrollmentInterceptor) InvalidateCache(agentId string) {
	a.cache.Del(agentId)
}

func (a *AgentEnrollmentInterceptor) Close() {
	a.cache.Close()
}

func (a *AgentEnrollmentInterceptor) requestEnrolTokenAgent(ctx context.Context, agentId string) (string, error) {
	value, found := a.cache.Get(agentId)
	if found {
		return value, nil
	}
	internalCtx := context.WithValue(ctx, models.ContextInternal, true)
	agent, err := a.agentService.GetAgent(internalCtx, agentId)
	if err != nil {
		return "", fmt.Errorf("unable to get agent: %w", err)
	}

	if agent.Status == models.AgentStatusEnrolling {
		return "", fmt.Errorf("agent is not enrolled")
	}

	a.cache.SetWithTTL(agentId, agent.EnrollmentToken, 1, 24*time.Hour)
	return agent.EnrollmentToken, nil
}
