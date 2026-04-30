package models

import "zhacked.me/oxyl/shared/pkg/models"

type ThresholdNotification struct {
	Identifier    string                  `json:"identifier"`
	AgentID       string                  `json:"agent_id"`
	TriggerReason models.NotificationType `json:"trigger_reason"`
	TriggerValue  string                  `json:"trigger_value"`
	Resolved      bool                    `json:"resolved"`
}
