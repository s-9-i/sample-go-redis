package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) error
	Decrement(ctx context.Context, key string) error
}

type cache struct {
	redisClient *redis.Client
}

func NewCache(c *redis.Client) Cache {
	return &cache{
		redisClient: c,
	}
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	v, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

func (c *cache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.redisClient.Set(ctx, key, value, ttl).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.redisClient.Del(ctx, key).Err()
}

func (c *cache) Increment(ctx context.Context, key string) error {
	return c.redisClient.Incr(ctx, key).Err()
}

func (c *cache) Decrement(ctx context.Context, key string) error {
	return c.redisClient.Decr(ctx, key).Err()
}
