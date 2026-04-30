package models

import "zhacked.me/oxyl/shared/pkg/models"

type ThresholdUpdate struct {
	CompanyId     string                  `json:"company_id"`
	ThresholdType models.NotificationType `json:"notification_type"`
	Threshold     float64                 `json:"threshold"`
}
