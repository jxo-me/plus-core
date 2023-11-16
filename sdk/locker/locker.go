package locker

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
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

// Mutex 返回分布式互斥锁对象
func (e *Locker) Mutex(ctx context.Context, key string, options ...redsync.Option) *redsync.Mutex {
	return e.locker.Mutex(ctx, e.getPrefixKey(key), options...)
}
