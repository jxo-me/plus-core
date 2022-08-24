package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/cache"
)

var insCache = Cache{}

type Cache struct {
	Redis  *GRedisOptions
	Memory interface{}
}

// CacheConfig cache配置
func CacheConfig() *Cache {
	return &insCache
}

// Setup 构造cache 顺序 redis > 其他 > memory
func (e *Cache) Setup(ctx context.Context) (storage.AdapterCache, error) {
	if e.Redis != nil {
		options, err := e.Redis.GetClientOptions()
		if err != nil {
			return nil, err
		}
		r, err := cache.NewGredis(GRedis().GetClient(), options)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return cache.NewMemory(), nil
}
