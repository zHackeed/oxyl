package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "zhacked.me/oxyl/protocol/v1"
	"zhacked.me/oxyl/protocol/v1/enrollment"
	"zhacked.me/oxyl/shared/pkg/datasource"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/utils"
	"zhacked.me/oxyl/shared/pkg/variables"
)

type EnrollmentService struct {
	v1.UnimplementedEnrollmentServiceServer

	messenger    *datasource.RedisConnection
	agentService *service.AgentService
	tokenService *service.TokenService
}

var _ v1.EnrollmentServiceServer = (*EnrollmentService)(nil)

func NewEnrollmentService(
	messenger *datasource.RedisConnection,
	agentService *service.AgentService,
	tokenService *service.TokenService,
) *EnrollmentService {
	return &EnrollmentService{
		messenger:    messenger,
		agentService: agentService,
		tokenService: tokenService,
	}
}

func (e *EnrollmentService) GetEnrollmentToken(ctx context.Context, req *enrollment.EnrollmentRequest) (*enrollment.EnrollmentResponse, error) {
	agentId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyAgent)
	if !found {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	internalCtx := context.WithValue(ctx, models.ContextInternal, true)
	agent, err := e.agentService.GetAgent(internalCtx, agentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "agent not found")
	}

	if agent.Status != models.AgentStatusEnrolling {
		slog.Info("agent is not enrolling", slog.String("agent_id", agentId), slog.String("status", string(agent.Status)))
		return nil, status.Error(codes.FailedPrecondition, "agent cannot enroll, it is already enrolled with current status")
	}

	// agent_id|os_name|cpu_model|total_memory|total_disk|update_timestamp_unix
	data := fmt.Sprintf("%s|%s|%s|%d|%d|%v", agentId,
		req.GetOsVariant(),
		req.GetCpuModel(), req.GetTotalMemory(),
		req.GetTotalDisk(), time.Now().Unix(),
	)
	signedToken := e.tokenService.GenerateSignedToken(data)

	// Convert the partitions to the internal model
	partitions := make([]*models.AgentPartition, 0, len(req.GetDiskPartitions()))
	for _, p := range req.GetDiskPartitions() {
		partitions = append(partitions, &models.AgentPartition{
			MountPoint: p.GetMountPoint(),
			TotalSize:  p.GetTotalSize(),
			Raid:       p.GetIsRaid(),
			RaidLevel:  int(p.GetRaidLevel()),
		})
	}

	if err := e.agentService.EnrichAgent(internalCtx, agentId,
		req.GetOsVariant(), req.GetCpuModel(),
		req.GetTotalMemory(), req.GetTotalDisk(),
		partitions, signedToken,
	); err != nil {
		slog.Error("unable to enrich agent", "error", err)
		return nil, status.Error(codes.Internal, "unable to enrich agent")
	}

	slog.Info("enrollment token generated", slog.String("agent_id", agentId))
	slog.Info("enrich data", slog.String("agent_id", agentId), slog.Any("data", data))

	if err := e.messenger.Publish(ctx, variables.RedisChannelAgentEnrollment, redisModels.AgentEnrollment{
		AgentId:      agentId,
		EnrollmentId: signedToken,
	}); err != nil {
		slog.Error("unable to publish agent enrollment event", "error", err)
	}

	return &enrollment.EnrollmentResponse{
		EnrollmentId: signedToken,
	}, nil
}
