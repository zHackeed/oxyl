package requests

import "zhacked.me/oxyl/shared/pkg/models"

type CreateAgentRequest struct {
	Holder       string `json:"holder"`
	DisplayName  string `json:"display_name"`
	RegisteredIP string `json:"registered_ip"`
}

type UpdateAgentStatusRequest struct {
	Status models.AgentStatus `json:"status"`
}
