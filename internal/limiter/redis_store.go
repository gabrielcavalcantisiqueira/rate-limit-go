package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStore{client: rdb}
}

func (r *RedisStore) Allow(id string, limit int, window time.Duration) bool {
	ctx := context.Background()
	key := fmt.Sprintf("limiter:reqs:%s", id)
	pipe := r.client.TxPipeline()
	now := time.Now().UnixNano()
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	})
	minScore := float64(time.Now().Add(-window).UnixNano())
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", minScore))
	count := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, window*2)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false
	}
	n, err := count.Result()
	if err != nil {
		return false
	}
	return n <= int64(limit)
}

func (r *RedisStore) IsLocked(id string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("limiter:lock:%s", id)
	untilStr, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	untilUnix, err := time.Parse(time.RFC3339Nano, untilStr)
	if err != nil {
		return false
	}
	return time.Now().Before(untilUnix)
}

func (r *RedisStore) SetLock(id string, until time.Time) {
	ctx := context.Background()
	key := fmt.Sprintf("limiter:lock:%s", id)
	r.client.Set(ctx, key, until.Format(time.RFC3339Nano), time.Until(until))
}

func (r *RedisStore) GetLockUntil(id string) time.Time {
	ctx := context.Background()
	key := fmt.Sprintf("limiter:lock:%s", id)
	untilStr, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339Nano, untilStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
