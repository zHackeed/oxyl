package requests

import "zhacked.me/oxyl/shared/pkg/models"

// todo: expand this maybe to have more information about the company.
type CreateCompanyRequest struct {
	DisplayName string `json:"display_name"`
}

// todo: validation
type AddMemberRequest struct {
	UserEmail  string                   `json:"user_id"`
	Permission models.CompanyPermission `json:"permission"`
}

type ModifyThresholdRequest struct {
	NotificationType models.NotificationType `json:"notification_type"`
	Threshold        int                     `json:"threshold"`
}
