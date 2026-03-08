package ratelimiter

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type WebhookRateLimiter interface {
	Allow(ctx context.Context, tenantID int64) (bool, error)
}

type RedisWebhookRateLimiter struct {
	client *goredis.Client
	limit  int64
}

func NewRedisWebhookRateLimiter(client *goredis.Client, limitPerMinute int) *RedisWebhookRateLimiter {
	if limitPerMinute <= 0 {
		limitPerMinute = 60
	}

	return &RedisWebhookRateLimiter{
		client: client,
		limit:  int64(limitPerMinute),
	}
}

func (r *RedisWebhookRateLimiter) Allow(ctx context.Context, tenantID int64) (bool, error) {
	minute := time.Now().UTC().Format("200601021504")
	key := fmt.Sprintf("rl:webhook:%d:%s", tenantID, minute)

	pipe := r.client.TxPipeline()
	countCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 2*time.Minute)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	return countCmd.Val() <= r.limit, nil
}
