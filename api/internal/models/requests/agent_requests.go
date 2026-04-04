package requests

import (
	"zhacked.me/oxyl/shared/pkg/models"
)

type AgentIdUri struct {
	AgentId string `uri:"id" validate:"required,alphanum"`
}

type CreateAgentRequest struct {
	Holder       string `json:"holder" validate:"required,alphanum"`
	DisplayName  string `json:"display_name" validate:"required"`
	RegisteredIP string `json:"registered_ip" validate:"required,ip"`
}

type UpdateAgentStatusRequest struct {
	Agent  string             `uri:"id" validate:"required,alphanum"`
	Status models.AgentStatus `json:"status" validate:"required,oneof=ACTIVE MAINTENANCE INACTIVE"`
}

type AgentLoginRequest struct {
	AgentId string `json:"agent_id" validate:"required,alphanum"`
}

type AuthenticationShutdownRequest struct {
	AgentId string `json:"agent_id" validate:"required,alphanum"`
}

type AgentMetricsRequest struct {
	AgentId  string `uri:"id" validate:"required,alphanum"`
	Interval string `uri:"interval" validate:"required,oneof=1m 5m 15m 1h 1d"` // todo: change to enum
}
