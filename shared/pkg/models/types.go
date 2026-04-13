package models

type WebhookType string

const (
	WebhookTypeWebhook WebhookType = "WEBHOOK"
	WebhookTypeDiscord WebhookType = "DISCORD"
	WebhookTypeSlack   WebhookType = "SLACK"
)

type NotificationType string

const (
	NotificationTypeCompanySettingUpdate NotificationType = "COMPANY_SETTING_UPDATE"
	NotificationTypeCompanyMemberUpdate  NotificationType = "COMPANY_MEMBER_UPDATE"

	NotificationTypeAgentStatusUpdate          NotificationType = "AGENT_STATUS_UPDATE"
	NotificationTypeAgentCpuUsageThreshold     NotificationType = "AGENT_CPU_USAGE_THRESHOLD"
	NotificationTypeAgentMemoryUsageThreshold  NotificationType = "AGENT_MEMORY_USAGE_THRESHOLD"
	NotificationTypeAgentDiskUsageThreshold    NotificationType = "AGENT_DISK_USAGE_THRESHOLD"
	NotificationTypeAgentDiskHealthThreshold   NotificationType = "AGENT_DISK_HEALTH_THRESHOLD"
	NotificationTypeAgentNetworkUsageThreshold NotificationType = "AGENT_NETWORK_USAGE_THRESHOLD"
)

func NotificationTypes() []NotificationType {
	return []NotificationType{
		NotificationTypeCompanySettingUpdate,
		NotificationTypeCompanyMemberUpdate,
		NotificationTypeAgentStatusUpdate,
		NotificationTypeAgentCpuUsageThreshold,
		NotificationTypeAgentMemoryUsageThreshold,
		NotificationTypeAgentDiskUsageThreshold,
		NotificationTypeAgentDiskHealthThreshold,
		NotificationTypeAgentNetworkUsageThreshold,
	}
}

const (
	AgentStatusActive      AgentStatus = "ACTIVE"
	AgentStatusEnrolling   AgentStatus = "ENROLLING"
	AgentStatusMaintenance AgentStatus = "MAINTENANCE"
	AgentStatusInactive    AgentStatus = "INACTIVE"
)

type AgentStatus string

const (
	TokenTypeAgent JWTTokenType = "AGENT"
	TokenTypeUser  JWTTokenType = "USER"
)

type JWTTokenType string

const (
	ContextKeyUser  ContextKey = "oxyl_user_identifier"
	ContextKeyAgent ContextKey = "oxyl_agent_identifier"
	ContextInternal ContextKey = "oxyl_internal_request"
)

type ContextKey string
