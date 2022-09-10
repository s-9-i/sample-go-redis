package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Latests interface {
	LPush(ctx context.Context, key string, val interface{}) (int64, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, count int64) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)
}

type latests struct {
	redisClient *redis.Client
}

func NewLatests(c *redis.Client) Latests {
	return &latests{
		redisClient: c,
	}
}

func (rl *latests) LPush(ctx context.Context, key string, val interface{}) (int64, error) {
	return rl.redisClient.LPush(ctx, key, val).Result()
}

func (rl *latests) RPop(ctx context.Context, key string) (string, error) {
	return rl.redisClient.RPop(ctx, key).Result()
}

func (rl *latests) LRange(ctx context.Context, key string, count int64) ([]string, error) {
	return rl.redisClient.LRange(ctx, key, 0, count-1).Result()
}

func (rl *latests) LLen(ctx context.Context, key string) (int64, error) {
	return rl.redisClient.LLen(ctx, key).Result()
}
