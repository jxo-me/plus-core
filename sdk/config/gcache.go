package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	cacheLib "github.com/jxo-me/plus-core/core/v2/cache"
	redisLib "github.com/jxo-me/plus-core/sdk/v2/cache/gredis"
	"github.com/jxo-me/plus-core/sdk/v2/cache/memory"
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
func (e *Cache) Setup(ctx context.Context, app app.IRuntime) (cacheLib.ICache, error) {
	redis := g.Redis(gredis.DefaultGroupName)
	if redis != nil {
		r, err := redisLib.NewGredis(redis)
		if err != nil {
			return nil, err
		}
		GRedis().SetClient(ctx, redis)
		return r, nil
	}
	options, err := e.Redis.GetClientOptions(ctx, app)
	if err != nil {
		return nil, err
	}
	redis, err = gredis.New(options)
	if err != nil {
		return nil, err
	}
	GRedis().SetClient(ctx, redis)
	r, err := redisLib.NewGredis(redis)
	if err != nil {
		glog.Warning(ctx, fmt.Sprintf("get redis cache options: %v error: %v", options, err))
		return memory.NewMemory(), nil
	}
	return r, nil
}
