package runtime

import (
	"github.com/jxo-me/plus-core/sdk/storage"
	"time"
)

const (
	intervalTenant = "/"
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
func (e Cache) Get(key string) (string, error) {
	return e.store.Get(e.prefix + intervalTenant + key)
}

// Set val in cache
func (e Cache) Set(key string, val interface{}, expire int) error {
	return e.store.Set(e.prefix+intervalTenant+key, val, expire)
}

// Del delete key in cache
func (e Cache) Del(key string) error {
	return e.store.Del(e.prefix + intervalTenant + key)
}

// HashGet get val in hashtable cache
func (e Cache) HashGet(hk, key string) (string, error) {
	return e.store.HashGet(hk, e.prefix+intervalTenant+key)
}

// HashDel delete one key:value pair in hashtable cache
func (e Cache) HashDel(hk, key string) error {
	return e.store.HashDel(hk, e.prefix+intervalTenant+key)
}

// Increase value
func (e Cache) Increase(key string) error {
	return e.store.Increase(e.prefix + intervalTenant + key)
}

func (e Cache) Decrease(key string) error {
	return e.store.Decrease(e.prefix + intervalTenant + key)
}

func (e Cache) Expire(key string, dur time.Duration) error {
	return e.store.Expire(e.prefix+intervalTenant+key, dur)
}
