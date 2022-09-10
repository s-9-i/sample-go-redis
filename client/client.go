package client

import (
	"github.com/go-redis/redis/v8"
)

type Options struct {
	Addr     string
	Password string
	DB       int
}

func NewClient(options *Options) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
}
