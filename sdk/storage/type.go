package storage

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"time"

	"github.com/bsm/redislock"
)

const (
	PrefixKey = "__host"
)

type AdapterCache interface {
	String() string
	Get(key string) (string, error)
	Set(key string, val interface{}, expire int) error
	Del(key string) error
	HashGet(hk, key string) (string, error)
	HashDel(hk, key string) error
	Increase(key string) error
	Decrease(key string) error
	Expire(key string, dur time.Duration) error
}

type AdapterGCache interface {
	String() string
	Get(ctx context.Context, key string) (*gvar.Var, error)
	Set(ctx context.Context, key string, val interface{}, expire int) error
	Del(ctx context.Context, key string) error
	HashGet(ctx context.Context, hk, key string) (*gvar.Var, error)
	HashDel(ctx context.Context, hk, key string) error
	Increase(ctx context.Context, key string) error
	Decrease(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, dur time.Duration) error
}

type AdapterQueue interface {
	String() string
	Append(ctx context.Context, message Messager) error
	Register(ctx context.Context, name string, f ConsumerFunc)
	Run(ctx context.Context)
	Shutdown(ctx context.Context)
}

type Messager interface {
	SetID(string)
	SetStream(string)
	SetValues(map[string]interface{})
	GetID() string
	GetStream() string
	GetValues() map[string]interface{}
	GetPrefix() string
	SetPrefix(string)
	SetErrorCount()
	GetErrorCount() int
}

type ConsumerFunc func(ctx context.Context, msg Messager) error

type AdapterLocker interface {
	String() string
	Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error)
}
