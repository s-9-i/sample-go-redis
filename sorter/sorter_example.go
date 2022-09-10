package main

import (
	"context"
	"fmt"
	"time"

	"sample-go-redis/client"
)

var s Sorter

const key = "key-sorter"

func main() {
	ctx := context.Background()
	cli := client.NewClient(&client.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		cli.Del(ctx, key)
		_ = cli.Close()
	}()

	s = NewSorter(cli)
	Sample(ctx)
}

func Sample(ctx context.Context) {
	s.Expire(ctx, key, time.Hour*24)

	s.Add(ctx, key, Members{
		{ID: "user1", Score: 100},
		{ID: "user2", Score: 200},
		{ID: "user3", Score: 300},
		{ID: "user4", Score: 300},
		{ID: "user5", Score: 500},
		{ID: "user6", Score: 600},
		{ID: "user7", Score: 600},
		{ID: "user8", Score: 600},
		{ID: "user9", Score: 900},
		{ID: "user10", Score: 1000},
	}...)

	rangeResults, _ := s.Range(ctx, key, 0, 9)
	for _, result := range rangeResults {
		fmt.Printf("【sorter】rangeResult = %#v\n", result)
	}

	rangeByScoreResults, _ := s.RangeByScore(ctx, key, 200, 800, 3, false)
	for _, result := range rangeByScoreResults {
		fmt.Printf("【sorter】rangeByScoreResults = %#v\n", result)
	}

	s.Delete(ctx, key, "user7", "user5")

	rangeByScoreRevResults, _ := s.RangeByScore(ctx, key, 200, 800, 3, true)
	for _, result := range rangeByScoreRevResults {
		fmt.Printf("【sorter】rangeByScoreRevResults = %#v\n", result)
	}
}
