package main

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Ranking interface {
	Add(ctx context.Context, key, id string, score int64, timestamp time.Time) error
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

func (rr *redisRanking) Add(ctx context.Context, key, id string, score int64, timestamp time.Time) error {
	pastTimestamp, err := rr.getPrevTimestamp(ctx, key, id)
	if err != nil {
		return err
	}
	oldMember := createRankingMember(pastTimestamp, id)

	newTimestamp := strconv.FormatInt(10_000_000_000_000-timestamp.UnixMilli(), 10)

	pipe := rr.redisClient.TxPipeline()
	if cmd := pipe.ZAdd(ctx, key, &redis.Z{Score: float64(score), Member: createRankingMember(newTimestamp, id)}); cmd.Err() != nil {
		return cmd.Err()
	}
	if cmd := pipe.Set(ctx, createTimestampKey(key, id), newTimestamp, 0); cmd.Err() != nil {
		return cmd.Err()
	}
	if oldMember != "" {
		if cmd := pipe.ZRem(ctx, key, oldMember); cmd.Err() != nil {
			return cmd.Err()
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (rr *redisRanking) Rank(ctx context.Context, key, id string) (int64, error) {
	pastTimestamp, err := rr.getPrevTimestamp(ctx, key, id)
	if err != nil {
		return 0, err
	}
	if pastTimestamp == "" {
		return 0, nil
	}
	val, err := rr.redisClient.ZRevRank(ctx, key, createRankingMember(pastTimestamp, id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return val + 1, nil
}

func (rr *redisRanking) RankByScore(ctx context.Context, key string, score int64) (int64, error) {
	count, err := rr.redisClient.ZCount(ctx, key, strconv.Itoa(int(score)), "+inf").Result()
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}

func (rr *redisRanking) Score(ctx context.Context, key, id string) (int64, error) {
	pastTimestamp, err := rr.getPrevTimestamp(ctx, key, id)
	if err != nil {
		return 0, err
	}
	if pastTimestamp == "" {
		return 0, nil
	}

	score, err := rr.redisClient.ZScore(ctx, key, createRankingMember(pastTimestamp, id)).Result()
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
		return nil, nil
	}

	results := make([]*RangeResult, 0, len(scores))
	rank := start + 1
	for _, result := range scores {
		results = append(results, &RangeResult{
			Rank:  rank,
			Score: int64(result.Score),
			ID:    convertToID(result.Member),
		})
		rank++
	}
	return results, nil
}

func (rr *redisRanking) getPrevTimestamp(ctx context.Context, key, id string) (string, error) {
	v, err := rr.redisClient.Get(ctx, createTimestampKey(key, id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

func createTimestampKey(key, id string) string {
	return key + "_" + id
}

func createRankingMember(timestamp, id string) string {
	if timestamp == "" {
		return ""
	}
	return timestamp + "_" + id
}

func convertToID(member interface{}) string {
	return strings.Split(member.(string), "_")[1]
}
