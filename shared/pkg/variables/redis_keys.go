package variables

const (
	redisKeyPrefix = "oxyl"

	//RedisKeyAgentJWT = redisKeyPrefix + ":agent:jwt"

	RedisTokenRevokedRedisKey = redisKeyPrefix + ":token:revoked"
)

type RedisKey string

const (
	RedisKeyHeartbeat       RedisKey = "agent:heartbeat:%s"
	RedisKeyThresholdActive RedisKey = "agent:thresholds:active:%s:%s"

	RedisChannelInvalidateUser RedisChannel = "user:invalidate"

	RedisChannelCompanyCreation        RedisChannel = "company:creation"
	RedisChannelCompanyAddedMember     RedisChannel = "company:member:add"
	RedisChannelCompanyRemovedMember   RedisChannel = "company:member:remove"
	RedisChannelCompanyThresholdUpdate RedisChannel = "company:threshold:update"

	RedisChannelCompanyWebhookCreate RedisChannel = "company:webhook:create"
	RedisChannelCompanyWebhookDelete RedisChannel = "company:webhook:delete"

	RedisChannelCompanyDeletion RedisChannel = "company:deletion"

	RedisChannelThresholdNotification RedisChannel = "agent:threshold:notification"

	RedisChannelAgentCreation RedisChannel = "company:agent:creation"
	RedisChannelAgentDeletion RedisChannel = "company:agent:deletion"

	RedisChannelAgentStateUpdate RedisChannel = "agent:state:update"

	RedisChannelAgentListening        RedisChannel = "agent:viewer:listening"
	RedisChannelAgentStoppedListening RedisChannel = "agent:viewer:deafen"

	RedisChannelAgentEnrollment RedisChannel = "agent:enrollment"
	RedisChannelAgentMetrics    RedisChannel = "agent:viewer:metric:append"
)

type RedisChannel string
