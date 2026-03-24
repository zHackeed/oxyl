package variables

const (
	redisKeyPrefix = "oxyl"

	//RedisKeyAgentJWT = redisKeyPrefix + ":agent:jwt"

	RedisTokenRevokedRedisKey = redisKeyPrefix + ":token:revoked"
)

type RedisKey string

const (
// RedisChannelAgentHeartbeat RedisChannel = "agent:heartbeat"
)

type RedisChannel string
