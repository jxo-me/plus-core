package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/cache"
)

type Cache struct {
	Redis  *GredisConnectOptions
	Memory interface{}
}

// CacheConfig cache配置
var CacheConfig = new(Cache)

// Setup 构造cache 顺序 redis > 其他 > memory
func (e Cache) Setup(ctx context.Context) (storage.AdapterGCache, error) {
	if e.Redis != nil {
		options, err := e.Redis.GetGredisOptions()
		if err != nil {
			return nil, err
		}
		r, err := cache.NewGredis(GetGredisClient(), options)
		if err != nil {
			return nil, err
		}
		if _redis == nil {
			_gredis = r.GetClient()
		}
		return r, nil
	}
	return cache.NewMemory(), nil
}
