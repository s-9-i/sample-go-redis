package main

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

var (
	ErrNotObtained = redislock.ErrNotObtained
	ErrLockNotHeld = redislock.ErrLockNotHeld
)

type Locker interface {
	TryLock(context.Context, string, time.Duration) (Lock, error)
}

type Lock interface {
	Release(context.Context) error
}

type locker struct {
	lockCli *redislock.Client
}

type lock struct {
	lock *redislock.Lock
}

func NewLocker(c *redis.Client) Locker {
	return &locker{
		lockCli: redislock.New(c),
	}
}

func NewLock(l *redislock.Lock) Lock {
	return &lock{
		lock: l,
	}
}

func (rl *locker) TryLock(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	l, err := rl.lockCli.Obtain(ctx, key, ttl, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(500*time.Millisecond), 10),
		Metadata:      "",
	})
	if err != nil {
		if err == redislock.ErrNotObtained {
			return nil, ErrNotObtained
		}
		return nil, err
	}
	return NewLock(l), nil
}

func (rl *lock) Release(ctx context.Context) error {
	if err := rl.lock.Release(ctx); err != nil {
		if err == redislock.ErrLockNotHeld {
			return ErrLockNotHeld
		}
		return err
	}
	return nil
}
