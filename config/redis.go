package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	DB *redis.Client
}

func NewRedisClient(redisUrl string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	return &RedisClient{DB: redis.NewClient(opts)}, nil
}

func (rclient *RedisClient) Set(key string, value any, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rclient.DB.Set(ctx, key, value, exp).Err()
}

func (rclient *RedisClient) Get(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rclient.DB.Get(ctx, key).Err()
}
