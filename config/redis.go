package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	DB *redis.Client
}

func newRedisClient(redisUrl string) (*redisClient, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	return &redisClient{DB: redis.NewClient(opts)}, nil
}

func (rclient *redisClient) Set(key string, value any, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rclient.DB.Set(ctx, key, value, exp).Err()
}

func (rclient *redisClient) Get(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rclient.DB.Get(ctx, key).Err()
}
