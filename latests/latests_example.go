package main

import (
	"context"
	"fmt"

	"sample-go-redis/client"
)

var l Latests

const key = "key-latests"

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

	l = NewLatests(cli)
	Sample(ctx)
}

func Sample(ctx context.Context) {
	l.LPush(ctx, key, "val1")
	l.LPush(ctx, key, "val2")

	vals, _ := l.LRange(ctx, key, 3)
	fmt.Printf("【latests】vals = %#v\n", vals)

	l.LPush(ctx, key, "val3")
	l.LPush(ctx, key, "val4")

	vals, _ = l.LRange(ctx, key, 3)
	fmt.Printf("【latests】vals = %#v\n", vals)
}
