package storage

import (
	"context"
	"time"

	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/variables"
)

type TokenStorage struct {
	conn *datasource.RedisConnection
}

func NewTokenStorage(persistence *datasource.RedisConnection) *TokenStorage {
	return &TokenStorage{
		conn: persistence,
	}
}

func (t *TokenStorage) RevokeToken(ctx context.Context, token *models.Token) error {
	return t.conn.HashSetIfNotExists(ctx, variables.RedisTokenRevokedRedisKey, token.Identifier, "1", 24*time.Hour)
}

func (t *TokenStorage) IsTokenRevoked(ctx context.Context, token *models.Token) (bool, error) {
	return t.conn.HashExists(ctx, variables.RedisTokenRevokedRedisKey, token.Identifier)
}
