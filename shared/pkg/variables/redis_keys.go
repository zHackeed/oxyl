package variables

const (
	redisKeyPrefix = "oxyl"

	//RedisKeyAgentJWT = redisKeyPrefix + ":agent:jwt"

	RedisTokenRevokedRedisKey = redisKeyPrefix + ":token:revoked"
)

type RedisKey string

const (
	// RedisChannelAgentHeartbeat RedisChannel = "agent:heartbeat"

	RedisChannelCompanyAddedMember     RedisChannel = "company:added_member"
	RedisChannelCompanyRemovedMember   RedisChannel = "company:removed_member"
	RedisChannelCompanyThresholdUpdate RedisChannel = "company:threshold_update"
	RedisChannelCompanyDeletion        RedisChannel = "company:deletion"

	RedisChannelAgentCreation RedisChannel = "agent:creation"
	// Todo: think about this.
	RedisChannelAgentEnrollment RedisChannel = "agent:enrollment"
	RedisChannelAgentUpdate     RedisChannel = "agent:update"
	RedisChannelAgentDeletion   RedisChannel = "agent:deletion"
	RedisChannelAgentHeartbeat  RedisChannel = "agent:heartbeat"
	RedisChannelAgentMetrics    RedisChannel = "agent:metrics"
)

type RedisChannel string
