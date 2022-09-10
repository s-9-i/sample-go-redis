package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Ranking interface {
	Add(ctx context.Context, key, id string, score int64) error
	Rank(ctx context.Context, key, id string) (int64, error)
	RankByScore(ctx context.Context, key string, score int64) (int64, error)
	Score(ctx context.Context, key, id string) (int64, error)
	Count(ctx context.Context, key string) (int64, error)
	Range(ctx context.Context, key string, start, end int64) ([]*RangeResult, error)
}

type redisRanking struct {
	redisClient *redis.Client
}

type RangeResult struct {
	Rank  int64
	Score int64
	ID    string
}

func NewRanking(c *redis.Client) Ranking {
	return &redisRanking{
		redisClient: c,
	}
}

func (rr *redisRanking) Add(ctx context.Context, key, id string, score int64) error {
	if cmd := rr.redisClient.ZAdd(ctx, key, &redis.Z{Score: float64(score), Member: id}); cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (rr *redisRanking) Rank(ctx context.Context, key, id string) (int64, error) {
	score, err := rr.Score(ctx, key, id)
	if err != nil {
		return 0, err
	}
	rank, err := rr.RankByScore(ctx, key, score)
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (rr *redisRanking) RankByScore(ctx context.Context, key string, score int64) (int64, error) {
	count, err := rr.redisClient.ZCount(ctx, key, strconv.Itoa(int(score)+1), "+inf").Result()
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}

func (rr *redisRanking) Score(ctx context.Context, key, id string) (int64, error) {
	score, err := rr.redisClient.ZScore(ctx, key, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return int64(score), nil
}

func (rr *redisRanking) Count(ctx context.Context, key string) (int64, error) {
	count, err := rr.redisClient.ZCount(ctx, key, "-inf", "+inf").Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (rr *redisRanking) Range(ctx context.Context, key string, start, end int64) ([]*RangeResult, error) {

	scores, err := rr.redisClient.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		return []*RangeResult{}, nil
	}

	results := make([]*RangeResult, 0, len(scores))
	baseRank := start + 1
	var rank int64
	var score int64
	for i, result := range scores {
		resultScore := int64(result.Score)
		if score != resultScore {
			score = resultScore
			rank = baseRank + int64(i)
		}
		results = append(results, &RangeResult{
			Rank:  rank,
			Score: score,
			ID:    result.Member.(string),
		})
	}
	return results, nil
}
