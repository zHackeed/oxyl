package messenger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/variables"
)

type PubSubRouter struct {
	listener  *datasource.RedisConnection
	handlers  map[variables.RedisChannel]handlerEntry
	inProcess sync.WaitGroup
}

type handlerEntry interface {
	getChannel() variables.RedisChannel
	handle(ctx context.Context, payload []byte) error
}

func NewPubSubRouter(redis *datasource.RedisConnection) *PubSubRouter {
	return new(PubSubRouter{
		listener: redis,
		handlers: make(map[variables.RedisChannel]handlerEntry),
	})
}

func RegisterHandler[T any](router *PubSubRouter, interceptor Interceptor[T]) {
	// We create a custom adapter for the interceptor so we can use it as a handler entry.
	// Thank you, golang limitations.
	adapter := &interceptorAdapter[T]{inner: interceptor}
	router.handlers[interceptor.GetChannel()] = adapter
}

func (r *PubSubRouter) Run(ctx context.Context) error {
	channels := make([]variables.RedisChannel, 0, len(r.handlers))
	for ch := range r.handlers {
		channels = append(channels, ch)
	}

	sub := r.listener.Subscribe(ctx, channels...)
	defer func() {
		// Close the channel and wait for all the handlers to finish.

		_ = sub.Close()
		r.inProcess.Wait()
	}()

	for {
		msg, err := sub.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("pubsub receive error: %w", err)
		}

		redisMsg, ok := msg.(*redis.Message)
		if !ok {
			continue
		}

		h, ok := r.handlers[variables.RedisChannel(redisMsg.Channel)]
		if !ok {
			continue
		}

		r.inProcess.Add(1)
		go r.handleMessage(context.WithoutCancel(ctx), h, redisMsg.Payload)
	}
}
func (r *PubSubRouter) handleMessage(ctx context.Context, h handlerEntry, payload string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in redis pubsub router", "error", r)
		}
		r.inProcess.Done()
	}()

	if err := h.handle(ctx, []byte(payload)); err != nil {
		slog.Error("error handling message", "channel", h.getChannel(), "error", err)
	}
}
