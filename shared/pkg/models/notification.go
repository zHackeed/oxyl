package models

import (
	"time"
)

type NotificationLog struct {
	Identifier    string           `json:"identifier"`
	Agent         string           `json:"agent"`
	TriggerReason NotificationType `json:"trigger_reason"`
	TriggerValue  string           `json:"trigger_value"`
	Ack           bool             `json:"ack"`
	Failed        bool             `json:"failed"`
	SentAt        time.Time        `json:"sent_at"`
}

type NotificationSetting struct {
	ID          string      `json:"id"`
	Holder      string      `json:"holder"`
	WebhookType WebhookType `json:"webhook_type"`
	Endpoint    string      `json:"endpoint"`
	Channel     *string     `json:"channel,omitempty"` // only for discord
}
