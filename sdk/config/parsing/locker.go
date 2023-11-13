package parsing

import (
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	lockerLib "github.com/jxo-me/plus-core/core/v2/locker"
	redis2 "github.com/jxo-me/plus-core/sdk/v2/locker/redis"
)

func ParseRedisLocker(cfg *gredis.Config) (lockerLib.ILocker, error) {
	redis := g.Redis(gredis.DefaultGroupName)
	return redis2.NewRedis(redis), nil
}
