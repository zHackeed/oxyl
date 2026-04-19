package models

import "zhacked.me/oxyl/shared/pkg/models"

type AgentCreation struct {
	CompanyId string             `json:"company_id"`
	AgentId   string             `json:"agent_id"`
	State     models.AgentStatus `json:"state"`
	// Might use the api to query? or poll periodically
	RegisteredIP string `json:"registered_ip"`
	DisplayName  string `json:"display_name"`
}

type AgentUpdate struct {
	CompanyHolder string             `json:"company_holder"`
	AgentId       string             `json:"agent_id"`
	Status        models.AgentStatus `json:"status"`
}

type AgentDelete struct {
	CompanyId string `json:"company_id"`
	AgentId   string `json:"agent_id"`
}

type AgentEnrollment struct {
	AgentId      string `json:"agent_id"`
	EnrollmentId string `json:"enrollment_id"`
}
