package responses

import (
	"time"

	"zhacked.me/oxyl/shared/pkg/models"
)

type CompanyMemberPermissionResponse struct {
	User        string   `json:"user_id"`
	Permissions []string `json:"permissions"`
}

type CompanyThresholdValueWrapper struct {
	ThresholdIdentifier models.NotificationType `json:"threshold_id"`
	Value               int                     `json:"value"`
}

type CompanyUserValueWrapper struct {
	User        models.UserResumed `json:"user"`
	Permissions []string           `json:"permissions"`
	CreatedAt   time.Time          `json:"created_at"`
}
