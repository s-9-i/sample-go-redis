package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"sample-go-redis/client"
)

var l Locker

const key = "key-lock"

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

	l = NewLocker(cli)
	SampleWithoutLock(ctx)
	fmt.Println()
	SampleWithLock(ctx)
}

func Hello(i int) {
	fmt.Println("【lock】Hello, Start: ", i)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("【lock】Hello, End: ", i)
}

func SampleWithoutLock(_ context.Context) {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			Hello(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func SampleWithLock(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			DoWithLock(ctx, func() {
				Hello(i)
			})
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func DoWithLock(ctx context.Context, f func()) {
	lockInfo, err := l.TryLock(ctx, key, 5*time.Second)
	if err != nil {
		if errors.Is(err, ErrNotObtained) {
			fmt.Printf("failed to obtain. err = %#v\n", err)
			return
		}
		fmt.Printf("failed to lock. err = %#v\n", err)
		return
	}
	defer func() {
		err := lockInfo.Release(ctx)
		if err != nil && !errors.Is(err, ErrLockNotHeld) {
			fmt.Printf("failed to release. err = %#v\n", err)
			return
		}
	}()
	f()
}
