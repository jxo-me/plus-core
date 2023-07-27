package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	lockerLib "github.com/jxo-me/plus-core/core/v2/locker"
	"github.com/jxo-me/plus-core/sdk/v2/locker/redis"
)

var insLocker = Locker{}

type Locker struct {
	Redis *GRedisOptions
}

func LockerConfig() *Locker {
	return &insLocker
}

// Empty 空设置
func (e *Locker) Empty() bool {
	return e.Redis == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e *Locker) Setup(ctx context.Context, s *Settings) (lockerLib.ILocker, error) {
	client := g.Redis(gredis.DefaultGroupName)
	if client != nil {
		return redis.NewRedis(client), nil
	}
	options, err := e.Redis.GetClientOptions(ctx, s)
	if err != nil {
		glog.Warning(ctx, fmt.Sprintf("get redis Locker options: %v error: %v", options, err))
		return nil, err
	}
	client, err = gredis.New(options)
	if err != nil {
		return nil, err
	}
	return redis.NewRedis(client), nil
}
