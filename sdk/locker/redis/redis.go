package redis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
)

// NewRedis 初始化locker
func NewRedis(c *redsync.Redsync) *Redis {
	return &Redis{
		client: c,
	}
}

type Redis struct {
	client *redsync.Redsync
}

func (Redis) String() string {
	return "redis"
}

func (r *Redis) Mutex(ctx context.Context, key string, options ...redsync.Option) *redsync.Mutex {
	mutex := r.client.NewMutex(key,
		options...,
	)
	return mutex
}
