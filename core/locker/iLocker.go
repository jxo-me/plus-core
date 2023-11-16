package locker

import (
	"context"
	"github.com/go-redsync/redsync/v4"
)

type ILocker interface {
	String() string
	Mutex(ctx context.Context, key string, options ...redsync.Option) *redsync.Mutex
}
