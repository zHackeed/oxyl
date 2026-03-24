package datasource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"zhacked.me/oxyl/shared/pkg/variables"
)

const NoExpiration = time.Duration(0)

var (
	ErrRedisKeyNotFound = errors.New("redis key not found")
)

type RedisConnection struct {
	conn *redis.Client
}

func NewRedisConnection() (*RedisConnection, error) {
	redisUrl, err := variables.GetValue(variables.RedisUri)
	if err != nil {
		return nil, fmt.Errorf("unable to create redis connection: %w", err)
	}

	ops, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to create redis connection: %w", err)
	}

	ops.MinIdleConns = 5
	ops.MaxIdleConns = 10
	ops.DialTimeout = 5 * time.Second
	ops.ReadTimeout = 15 * time.Second
	ops.WriteTimeout = 15 * time.Second

	conn := redis.NewClient(ops)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// is redis alive and does the params work?
	if err := conn.Ping(timeoutCtx).Err(); err != nil {
		return nil, fmt.Errorf("unable to create redis connection: %w", err)
	}

	return &RedisConnection{
		conn: conn,
	}, nil
}

func (rc *RedisConnection) Get(ctx context.Context, key string) (string, error) {
	data, err := rc.conn.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrRedisKeyNotFound // might change.
		}

		return "", fmt.Errorf("unable to get redis key %q: %w", key, err)
	}

	return data, nil
}

func (rc *RedisConnection) Set(ctx context.Context, key variables.RedisKey, value string, expiration time.Duration) error {
	if len(value) == 0 {
		return errors.New("value is empty")
	}

	return rc.conn.Set(ctx, string(key), value, expiration).Err()
}

func (rc *RedisConnection) Exists(ctx context.Context, key variables.RedisKey) (bool, error) {
	exists, err := rc.conn.Exists(ctx, string(key)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("unable to check if redis key %q exists: %w", string(key), err)
	}

	return exists > 0, nil
}

func (rc *RedisConnection) Del(ctx context.Context, key variables.RedisKey) error {
	return rc.conn.Del(ctx, string(key)).Err()
}

func (rc *RedisConnection) HashSetIfNotExists(ctx context.Context, key variables.RedisKey, field string, value any, expiration time.Duration) error {
	return rc.conn.HSetNX(ctx, string(key), field, value).Err()
}

func (rc *RedisConnection) HashExists(ctx context.Context, key variables.RedisKey, field string) (bool, error) {
	return rc.conn.HExists(ctx, string(key), field).Result()
}

func (rc *RedisConnection) HashGet(ctx context.Context, key variables.RedisKey, field string) (string, error) {
	return rc.conn.HGet(ctx, string(key), field).Result()
}

func (rc *RedisConnection) Publish(ctx context.Context, channel variables.RedisChannel, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("unable to publish to redis channel %q: %w", channel, err)
	}

	return rc.conn.Publish(ctx, string(channel), data).Err()
}

func (rc *RedisConnection) Subscribe(ctx context.Context, channel ...variables.RedisChannel) *redis.PubSub {
	parsedChannels := make([]string, len(channel))
	for i, c := range channel {
		parsedChannels[i] = string(c)
	}

	return rc.conn.Subscribe(ctx, parsedChannels...)
}

func (rc *RedisConnection) Close() error {
	return rc.conn.Close()
}

/*
// These are typed values function calls so we can avoid all the data passing and such that we would need to do with the raw redis connection.
// This is not the most efficient way, but it is the most convenient for now.
// The plans might be to use this for models, but the structure is not yet defined at this point fully.

func Set[T any](ctx context.Context, conn *RedisConnection, key string, value T, expiration ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("unable to set redis key: %w", err)
	}

	if len(expiration) > 0 {
		return conn.set(ctx, key, data, expiration[0]).Err()
	}

	return conn.set(ctx, key, data, 0).Err()
}

func Pool[T any](ctx context.Context, conn *RedisConnection, key string) (T, error) {
	var value T
	data, err := conn.get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return value, fmt.Errorf("redis key %q not found", key)
		}

		return value, fmt.Errorf("unable to get redis key %q: %w", key, err)
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("unable to get redis key %q: %w", key, err)
	}

	return value, nil
}
*/
