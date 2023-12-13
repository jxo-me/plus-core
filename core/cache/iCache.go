package cache

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"time"
)

type ICache interface {
	String() string
	Get(ctx context.Context, key string) (*gvar.Var, error)
	Set(ctx context.Context, key string, val interface{}, expire int) error
	Del(ctx context.Context, key string) error
	HashGet(ctx context.Context, hk, key string) (*gvar.Var, error)
	HashDel(ctx context.Context, hk, key string) error
	HashSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error)
	HashLen(ctx context.Context, key string) (int64, error)
	HashGetAll(ctx context.Context, key string) (*gvar.Var, error)
	HashMSet(ctx context.Context, key string, fields map[string]interface{}) error
	HashMGet(ctx context.Context, key string, fields ...string) (gvar.Vars, error)
	Increase(ctx context.Context, key string) (int64, error)
	Decrease(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, dur time.Duration) error
}
