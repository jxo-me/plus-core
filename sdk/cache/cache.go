package cache

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/jxo-me/plus-core/core/v2/cache"
	"time"
)

// NewCache 创建对应上下文缓存
func NewCache(prefix string, store cache.ICache) cache.ICache {
	return &Cache{
		prefix: prefix,
		store:  store,
	}
}

type Cache struct {
	prefix string
	store  cache.ICache
}

// String string输出
func (e *Cache) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

func (e *Cache) getPrefixKey(key string) string {
	if e.prefix != "" {
		return fmt.Sprintf("%s:%s", e.prefix, key)
	}
	return key
}

// SetPrefix 设置前缀
func (e *Cache) SetPrefix(prefix string) {
	e.prefix = prefix
}

// Connect 初始化
func (e *Cache) Connect() error {
	return nil
	//return e.store.Connect()
}

// Get val in cache
func (e *Cache) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return e.store.Get(ctx, e.getPrefixKey(key))
}

// Set val in cache
func (e *Cache) Set(ctx context.Context, key string, val interface{}, expire int) error {
	return e.store.Set(ctx, e.getPrefixKey(key), val, expire)
}

// Del delete key in cache
func (e *Cache) Del(ctx context.Context, key string) error {
	return e.store.Del(ctx, e.getPrefixKey(key))
}

// HashGet get val in hashtable cache
func (e *Cache) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	return e.store.HashGet(ctx, hk, e.getPrefixKey(key))
}

func (e *Cache) HashSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error) {
	return e.store.HashSet(ctx, key, fields)
}

func (e *Cache) HashMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	return e.store.HashMSet(ctx, key, fields)
}

func (e *Cache) HashMGet(ctx context.Context, key string, fields ...string) (gvar.Vars, error) {
	return e.store.HashMGet(ctx, key, fields...)
}

func (e *Cache) HashVals(ctx context.Context, key string) (gvar.Vars, error) {
	return e.store.HashVals(ctx, key)
}

func (e *Cache) HashLen(ctx context.Context, key string) (int64, error) {
	return e.store.HashLen(ctx, key)
}

func (e *Cache) HashGetAll(ctx context.Context, key string) (*gvar.Var, error) {
	return e.store.HashGetAll(ctx, key)
}

// HashDel delete one key:value pair in hashtable cache
func (e *Cache) HashDel(ctx context.Context, hk, key string) error {
	return e.store.HashDel(ctx, hk, e.getPrefixKey(key))
}

// Increase value
func (e *Cache) Increase(ctx context.Context, key string) (int64, error) {
	return e.store.Increase(ctx, e.getPrefixKey(key))
}

func (e *Cache) Decrease(ctx context.Context, key string) (int64, error) {
	return e.store.Decrease(ctx, e.getPrefixKey(key))
}

func (e *Cache) Expire(ctx context.Context, key string, dur time.Duration) error {
	return e.store.Expire(ctx, e.getPrefixKey(key), dur)
}
