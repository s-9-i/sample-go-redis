package main

import (
	"context"
	"fmt"
	"time"

	"sample-go-redis/client"
)

var r Ranking

const key = "key-ranking-timestamp"

func main() {
	ctx := context.Background()
	cli := client.NewClient(&client.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		cli.Del(ctx, cli.Keys(ctx, key+"*").Val()...)
		_ = cli.Close()
	}()

	r = NewRanking(cli)
	Sample(ctx)
}

func Sample(ctx context.Context) {
	now := time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)
	for _, user := range []struct {
		ID        string
		Score     int64
		TimeStamp time.Time
	}{
		{ID: "user1", Score: 100, TimeStamp: now.AddDate(0, 0, -10)},
		{ID: "user2", Score: 200, TimeStamp: now.AddDate(0, 0, -9)},
		{ID: "user3", Score: 300, TimeStamp: now.AddDate(0, 0, -8)},
		{ID: "user4", Score: 300, TimeStamp: now.AddDate(0, 0, -7)},
		{ID: "user5", Score: 500, TimeStamp: now.AddDate(0, 0, -6)},
		{ID: "user6", Score: 600, TimeStamp: now.AddDate(0, 0, -5)},
		{ID: "user7", Score: 600, TimeStamp: now.AddDate(0, 0, -4)},
		{ID: "user8", Score: 600, TimeStamp: now.AddDate(0, 0, -3)},
		{ID: "user9", Score: 900, TimeStamp: now.AddDate(0, 0, -2)},
		{ID: "user10", Score: 1000, TimeStamp: now.AddDate(0, 0, -1)},
	} {
		r.Add(ctx, key, user.ID, user.Score, user.TimeStamp)
	}

	rangeResults, _ := r.Range(ctx, key, 0, 9)
	for _, result := range rangeResults {
		fmt.Printf("【ranking】rangeResult = %#v\n", result)
	}

	rankResult, _ := r.Rank(ctx, key, "user6")
	fmt.Printf("【ranking】rankResult(user6) = %#v\n", rankResult)
	rankResult, _ = r.Rank(ctx, key, "user7")
	fmt.Printf("【ranking】rankResult(user7) = %#v\n", rankResult)
	rankResult, _ = r.Rank(ctx, key, "user8")
	fmt.Printf("【ranking】rankResult(user8) = %#v\n", rankResult)

	rankByScoreResult, _ := r.RankByScore(ctx, key, 600)
	fmt.Printf("【ranking】rankByScoreResult = %#v\n", rankByScoreResult)

	scoreResult, _ := r.Score(ctx, key, "user7")
	fmt.Printf("【ranking】scoreResult = %#v\n", scoreResult)

	countResult, _ := r.Count(ctx, key)
	fmt.Printf("【ranking】countResult = %#v\n", countResult)
}
