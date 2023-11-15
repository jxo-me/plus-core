package parsing

import (
	lockerLib "github.com/jxo-me/plus-core/core/v2/locker"
	"github.com/jxo-me/plus-core/sdk/v2/config"
	redisLocker "github.com/jxo-me/plus-core/sdk/v2/locker/redis"
	redisLib "github.com/redis/go-redis/v9"
)

func ParseRedisLocker(cfg *config.RedisOptions) (lockerLib.ILocker, error) {
	opt := redisLib.Options{
		Network:  "tcp",
		Addr:     cfg.Address,
		Username: cfg.User,
		Password: cfg.Pass,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	}

	return redisLocker.NewRedis(redisLib.NewClient(&opt)), nil
}
