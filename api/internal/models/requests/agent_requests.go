package requests

import "zhacked.me/oxyl/shared/pkg/models"

type AgentIdUri struct {
	AgentId string `uri:"id" validate:"required,alphanumeric"`
}

type CreateAgentRequest struct {
	Holder       string `json:"holder" validate:"required,alphanumeric"`
	DisplayName  string `json:"display_name" validate:"required"`
	RegisteredIP string `json:"registered_ip" validate:"required,ip"`
}

type UpdateAgentStatusRequest struct {
	Agent  string             `uri:"id" validate:"required,alphanumeric"`
	Status models.AgentStatus `json:"status" validate:"required,oneof=ACTIVE MAINTENANCE INACTIVE"`
}
