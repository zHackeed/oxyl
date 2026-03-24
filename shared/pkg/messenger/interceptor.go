package messenger

import (
	"context"
	"encoding/json"
	"fmt"

	"zhacked.me/oxyl/shared/pkg/variables"
)

type Interceptor[T any] interface {
	GetChannel() variables.RedisChannel
	Intercept(ctx context.Context, msg T) error
}

// So golang is not the best at abstraction and generics. So if we need to do this, we need to adapt the interface to the desired type.
// So if we can consume the data, it becomes a lot easier to work with as we are not duplicating the whole serialization and deserialization logic.

// Via this way we can just have a structs tha implements the interface and have the logic to handle the redis message with the desired type.

type interceptorAdapter[T any] struct {
	inner Interceptor[T]
}

func (a *interceptorAdapter[T]) getChannel() variables.RedisChannel {
	return a.inner.GetChannel()
}

func (a *interceptorAdapter[T]) handle(ctx context.Context, payload []byte) error {
	var msg T
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("unmarshal failed on channel %q: %w", a.getChannel(), err)
	}

	return a.inner.Intercept(ctx, msg)
}
