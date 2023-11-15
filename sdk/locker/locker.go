package locker

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/jxo-me/plus-core/core/v2/locker"
)

// NewLocker 创建对应上下文分布式锁
func NewLocker(prefix string, locker locker.ILocker) locker.ILocker {
	return &Locker{
		prefix: prefix,
		locker: locker,
	}
}

type Locker struct {
	prefix string
	locker locker.ILocker
}

func (e *Locker) String() string {
	return e.locker.String()
}

func (e *Locker) getPrefixKey(key string) string {
	return fmt.Sprintf("%s:%s", e.prefix, key)
}

// Lock 返回分布式锁对象
func (e *Locker) Lock(ctx context.Context, key string, ttl int64, options *redislock.Options) (*redislock.Lock, error) {
	return e.locker.Lock(ctx, e.getPrefixKey(key), ttl, options)
}
