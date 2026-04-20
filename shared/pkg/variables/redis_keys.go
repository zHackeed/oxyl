package variables

const (
	redisKeyPrefix = "oxyl"

	//RedisKeyAgentJWT = redisKeyPrefix + ":agent:jwt"

	RedisTokenRevokedRedisKey = redisKeyPrefix + ":token:revoked"
)

type RedisKey string

const (
	RedisChannelInvalidateUser RedisChannel = "user:invalidate"

	RedisChannelCompanyAddedMember     RedisChannel = "company:added_member"
	RedisChannelCompanyRemovedMember   RedisChannel = "company:removed_member"
	RedisChannelCompanyThresholdUpdate RedisChannel = "company:threshold_update"
	RedisChannelCompanyDeletion        RedisChannel = "company:deletion"

	RedisChannelAgentCreation   RedisChannel = "company:agent:creation"
	RedisChannelAgentEnrollment RedisChannel = "agent:enrollment"
	RedisChannelAgentUpdate     RedisChannel = "company:agent:update"
	RedisChannelAgentDeletion   RedisChannel = "company:agent:deletion"
	RedisChannelAgentHeartbeat  RedisChannel = "agent:heartbeat"
	RedisChannelAgentMetrics    RedisChannel = "agent:metrics"
)

type RedisChannel string
