package main

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Sorter interface {
	Add(ctx context.Context, key string, members ...*Member) error
	Range(ctx context.Context, key string, start, end int64) (Members, error)
	RangeByScore(ctx context.Context, key string, min, max, count int64, rev bool) (Members, error)
	Delete(ctx context.Context, key string, ids ...interface{}) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

type Member struct {
	Score float64
	ID    string
}

type Members []*Member

type redisSorter struct {
	redisClient *redis.Client
}

func NewSorter(c *redis.Client) Sorter {
	return &redisSorter{
		redisClient: c,
	}
}

func (rs *redisSorter) Add(ctx context.Context, key string, members ...*Member) error {
	zs := make([]*redis.Z, 0, len(members))
	for _, member := range members {
		zs = append(zs, &redis.Z{
			Score:  member.Score,
			Member: member.ID,
		})
	}

	zaddCmd := rs.redisClient.ZAdd(ctx, key, zs...)
	if zaddCmd.Err() != nil {
		return zaddCmd.Err()
	}

	return nil
}

func (rs *redisSorter) Range(ctx context.Context, key string, start, end int64) (Members, error) {
	scores, err := rs.redisClient.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		return Members{}, nil
	}

	results := make(Members, 0, len(scores))
	rank := start + 1
	for _, result := range scores {
		results = append(results, &Member{
			Score: result.Score,
			ID:    result.Member.(string),
		})
		rank++
	}
	return results, nil
}

func (rs *redisSorter) RangeByScore(ctx context.Context, key string, min, max, count int64, rev bool) (Members, error) {
	rangeBy := &redis.ZRangeBy{
		Min:    strconv.Itoa(int(min)),
		Max:    strconv.Itoa(int(max)),
		Offset: 0,
		Count:  count,
	}
	var cmd *redis.ZSliceCmd
	if rev {
		cmd = rs.redisClient.ZRevRangeByScoreWithScores(ctx, key, rangeBy)
	} else {
		cmd = rs.redisClient.ZRangeByScoreWithScores(ctx, key, rangeBy)
	}
	scores, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	results := make(Members, 0, len(scores))
	for _, result := range scores {
		results = append(results, &Member{
			Score: result.Score,
			ID:    result.Member.(string),
		})
	}
	return results, nil
}

func (rs *redisSorter) Delete(ctx context.Context, key string, ids ...interface{}) error {
	return rs.redisClient.ZRem(ctx, key, ids...).Err()
}

func (rs *redisSorter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rs.redisClient.Expire(ctx, key, expiration).Err()
}
