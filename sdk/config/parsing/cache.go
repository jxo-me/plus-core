package parsing

import (
	"github.com/gogf/gf/v2/database/gredis"
	cacheLib "github.com/jxo-me/plus-core/core/v2/cache"
	redisLib "github.com/jxo-me/plus-core/sdk/v2/cache/gredis"
	"github.com/jxo-me/plus-core/sdk/v2/cache/memory"
	redis2 "github.com/jxo-me/plus-core/sdk/v2/cache/redis"
	"github.com/jxo-me/plus-core/sdk/v2/config"
	redisLib2 "github.com/redis/go-redis/v9"
)

func ParseRedisCache(cfg *config.RedisOptions) (cacheLib.ICache, error) {
	opt := redisLib2.Options{
		Network:  "tcp",
		Addr:     cfg.Address,
		Username: cfg.User,
		Password: cfg.Pass,
		DB:       cfg.Db,
	}

	return redis2.NewRedis(nil, &opt)
}

func ParseGredisCache(cfg *config.RedisOptions) (cacheLib.ICache, error) {
	redis, err := gredis.New(&cfg.Config)
	if err != nil {
		return nil, err
	}
	r, err := redisLib.NewGredis(redis)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ParseMemoryCache(cfg *config.MemoryOptions) (cacheLib.ICache, error) {
	return memory.NewMemory(), nil
}
