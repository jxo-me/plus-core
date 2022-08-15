package runtime

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/jxo-me/plus-core/sdk/storage"
	"time"
)

const (
	interval = "/"
)

// NewCache 创建对应上下文缓存
func NewCache(prefix string, store storage.AdapterCache, tokenStoreKey string) storage.AdapterCache {
	if tokenStoreKey == "" {
		tokenStoreKey = "token_store_key"
	}
	return &Cache{
		prefix:        prefix,
		store:         store,
		tokenStoreKey: tokenStoreKey,
	}
}

type Cache struct {
	prefix        string
	store         storage.AdapterCache
	tokenStoreKey string
}

// String string输出
func (e *Cache) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// SetPrefix 设置前缀
func (e *Cache) SetPrefix(prefix string) {
	e.prefix = prefix
}

// Connect 初始化
func (e Cache) Connect() error {
	return nil
	//return e.store.Connect()
}

// Get val in cache
func (e Cache) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return e.store.Get(ctx, e.prefix+interval+key)
}

// Set val in cache
func (e Cache) Set(ctx context.Context, key string, val interface{}, expire int) error {
	return e.store.Set(ctx, e.prefix+interval+key, val, expire)
}

// Del delete key in cache
func (e Cache) Del(ctx context.Context, key string) error {
	return e.store.Del(ctx, e.prefix+interval+key)
}

// HashGet get val in hashtable cache
func (e Cache) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	return e.store.HashGet(ctx, hk, e.prefix+interval+key)
}

// HashDel delete one key:value pair in hashtable cache
func (e Cache) HashDel(ctx context.Context, hk, key string) error {
	return e.store.HashDel(ctx, hk, e.prefix+interval+key)
}

// Increase value
func (e Cache) Increase(ctx context.Context, key string) error {
	return e.store.Increase(ctx, e.prefix+interval+key)
}

func (e Cache) Decrease(ctx context.Context, key string) error {
	return e.store.Decrease(ctx, e.prefix+interval+key)
}

func (e Cache) Expire(ctx context.Context, key string, dur time.Duration) error {
	return e.store.Expire(ctx, e.prefix+interval+key, dur)
}
