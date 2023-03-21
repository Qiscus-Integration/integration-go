package cache

import (
	"context"
	"time"

	rd "github.com/go-redis/redis"
)

type redis struct {
	client *rd.Client
}

// NewCacheRedis returns a new instance of Redis implementation of CacheRepository.
func NewCacheRedis(client *rd.Client) *redis {
	return &redis{client}
}

func (r *redis) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return r.client.Set(key, value, exp).Err()
}

func (r *redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *redis) Del(ctx context.Context, key string) error {
	return r.client.Del(key).Err()
}
