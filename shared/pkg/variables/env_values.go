package variables

import (
	"errors"
	"fmt"
	"os"
)

const (
	RedisUri    EnvKey = "REDIS_URI"
	TigerdbHost EnvKey = "TIGERDB_HOST"
	TigerdbPort EnvKey = "TIGERDB_PORT"
	TigerdbUser EnvKey = "TIGERDB_USER"
	TigerdbPass EnvKey = "TIGERDB_PASS"
	TigerdbDb   EnvKey = "TIGERDB_DB"
)

type EnvKey string

var notFound = errors.New("environment variable not found")

func GetValue(key EnvKey) (string, error) {
	variable := os.Getenv(string(key))
	if variable == "" {
		return "", notFound
	}

	return variable, nil
}

func GetValueAggregate(keys ...EnvKey) (map[EnvKey]string, error) {
	values := make(map[EnvKey]string, len(keys))

	for _, key := range keys {
		value, err := GetValue(key)
		if err != nil {
			return nil, fmt.Errorf("missing required key %q: %w", key, err)
		}
		values[key] = value
	}

	return values, nil
}
