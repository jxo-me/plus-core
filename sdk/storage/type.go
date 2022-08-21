package storage

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
	"time"

	"github.com/jxo-me/redislock"
)

const (
	PrefixKey = "__host"
)

type AdapterCache interface {
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
	Publish(ctx context.Context, message Messager, optionFuncs ...func(*queue.PublishOptions)) error
	Consumer(ctx context.Context, name string, f ConsumerFunc, optionFuncs ...func(*queue.ConsumeOptions))
	Run(ctx context.Context)
	Shutdown(ctx context.Context)
}

type Messager interface {
	SetId(string)
	GetId() string
	SetRoutingKey(string)
	GetRoutingKey() string
	SetValues(map[string]interface{})
	GetValues() map[string]interface{}
	GetPrefix() string
	SetPrefix(string)
	SetErrorCount()
	GetErrorCount() int
}

type ConsumerFunc func(ctx context.Context, msg Messager) error

type AdapterLocker interface {
	String() string
	Lock(key string, ttl int64, options ...redislock.Option) (*redislock.Mutex, error)
}
