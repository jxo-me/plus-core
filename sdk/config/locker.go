package config

import (
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/locker"
)

var LockerConfig = new(Locker)

type Locker struct {
	Redis *GredisConnectOptions
}

// Empty 空设置
func (e Locker) Empty() bool {
	return e.Redis == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e Locker) Setup() (storage.AdapterLocker, error) {
	if e.Redis != nil {
		client := GetGredisClient()
		if client == nil {
			options, err := e.Redis.GetGredisOptions()
			if err != nil {
				return nil, err
			}
			client, err = gredis.New(options)
			if err != nil {
				return nil, err
			}
			_gredis = client
		}
		return locker.NewRedis(client), nil
	}
	return nil, nil
}
