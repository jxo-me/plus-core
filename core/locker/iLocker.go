package locker

import (
	"context"
	"github.com/bsm/redislock"
)

type ILocker interface {
	String() string
	Lock(ctx context.Context, key string, ttl int64, options *redislock.Options) (*redislock.Lock, error)
}
