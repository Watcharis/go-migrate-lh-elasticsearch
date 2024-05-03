package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	redisClient *redis.Client
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
}

func NewCache(redisClient *redis.Client) Cache {
	return &cache{
		redisClient: redisClient,
	}
}

func (r *cache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
