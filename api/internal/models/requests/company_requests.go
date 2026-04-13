package requests

import "zhacked.me/oxyl/shared/pkg/models"

// todo: expand this maybe to have more information about the company.
type CreateCompanyRequest struct {
	DisplayName string `json:"display_name"`
}

// todo: validation
type AddMemberRequest struct {
	CompanyId  string                   `uri:"id" validate:"required,alphanum"`
	UserEmail  string                   `json:"user_email" validate:"required,email"`
	Permission models.CompanyPermission `json:"permission" validate:"required,oneof=1 2 4 8 16 32 63 999"`
}

type RemoveMemberRequest struct {
	CompanyId string `uri:"company_id"`
	UserID    string `uri:"user_id"`
}

type ModifyThresholdRequest struct {
	CompanyId        string                  `uri:"id" validate:"required,alphanum"`
	NotificationType models.NotificationType `json:"notification_type" validate:"required,oneof=COMPANY_SETTING_UPDATE COMPANY_MEMBER_UPDATE AGENT_STATUS_UPDATE AGENT_CPU_USAGE_THRESHOLD AGENT_MEMORY_USAGE_THRESHOLD AGENT_DISK_USAGE_THRESHOLD AGENT_DISK_HEALTH_THRESHOLD AGENT_NETWORK_USAGE_THRESHOLD"`
	Threshold        int                     `json:"threshold" validate:"required,numeric"`
}

type CompanyIdUri struct {
	CompanyId string `uri:"id" uri:"company_id" validate:"required,alphanum"`
}
