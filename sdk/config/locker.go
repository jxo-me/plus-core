package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/locker"
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
func (e *Locker) Setup(ctx context.Context, s *Settings) (storage.AdapterLocker, error) {
	client := g.Redis(gredis.DefaultGroupName)
	if client != nil {
		return locker.NewRedis(client), nil
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
	return locker.NewRedis(client), nil
}
