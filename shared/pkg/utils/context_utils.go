package utils

import (
	"context"

	"zhacked.me/oxyl/shared/pkg/models"
)

func GetValueFromContext[T any](ctx context.Context, key models.ContextKey) (T, bool) {
	value, ok := ctx.Value(key).(T)

	if !ok {
		var zero T
		return zero, false
	}

	return value, ok
}
