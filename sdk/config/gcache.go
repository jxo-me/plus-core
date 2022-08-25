package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
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
func (e *Cache) Setup(ctx context.Context, s *Settings) (storage.AdapterCache, error) {
	redis := g.Redis(gredis.DefaultGroupName)
	if redis != nil {
		r, err := cache.NewGredis(redis, nil)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	options, err := e.Redis.GetClientOptions(ctx, s)
	if err != nil {
		return nil, err
	}
	r, err := cache.NewGredis(redis, options)
	if err != nil {
		glog.Warning(ctx, fmt.Sprintf("get redis cache options: %v error: %v", options, err))
		return cache.NewMemory(), nil
	}
	return r, nil
}
