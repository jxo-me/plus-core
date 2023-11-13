package parsing

import (
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	cacheLib "github.com/jxo-me/plus-core/core/v2/cache"
	redisLib "github.com/jxo-me/plus-core/sdk/v2/cache/gredis"
	"github.com/jxo-me/plus-core/sdk/v2/cache/memory"
	redis2 "github.com/jxo-me/plus-core/sdk/v2/cache/redis"
	"github.com/jxo-me/plus-core/sdk/v2/config"
	redisLib2 "github.com/redis/go-redis/v9"
)

func ParseRedis(cfg *gredis.Config) (cacheLib.ICache, error) {
	opt := redisLib2.Options{
		Network:  "tcp",
		Addr:     cfg.Address,
		Username: cfg.User,
		Password: cfg.Pass,
		DB:       cfg.Db,
	}

	return redis2.NewRedis(nil, &opt)
}

func ParseGredis(cfg *gredis.Config) (cacheLib.ICache, error) {
	redis := g.Redis(gredis.DefaultGroupName)
	r, err := redisLib.NewGredis(redis)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ParseMemory(cfg *config.MemoryOptions) (cacheLib.ICache, error) {
	return memory.NewMemory(), nil
}
