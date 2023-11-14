package config

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"sync"
)

var (
	global = &Config{
		Database: make(map[string]*gdb.ConfigNode),
		Redis:    map[string]*gredis.Config{},
		Server:   &ghttp.ServerConfig{},
		Settings: &SettingOptions{
			Queue: &QueueGroups{
				Rocketmq: map[string]*RocketmqOptions{},
				Rabbitmq: map[string]*RabbitmqOptions{},
			},
			Auth:        map[string]*JwtAuth{},
			FailedLimit: map[string]*FailedLimitOptions{},
			Uploads:     &UploadGroups{Tus: map[string]*TusOptions{}},
		},
		Bot: &BotGroups{},
	}
	globalMux sync.RWMutex
)

func Global() *Config {
	globalMux.RLock()
	defer globalMux.RUnlock()

	cfg := &Config{
		Database: make(map[string]*gdb.ConfigNode),
		Redis:    map[string]*gredis.Config{},
		Server:   &ghttp.ServerConfig{},
		Settings: &SettingOptions{
			Queue: &QueueGroups{
				Rocketmq: map[string]*RocketmqOptions{},
				Rabbitmq: map[string]*RabbitmqOptions{},
			},
			Auth:        map[string]*JwtAuth{},
			FailedLimit: map[string]*FailedLimitOptions{},
			Uploads:     &UploadGroups{Tus: map[string]*TusOptions{}},
		},
		Bot: &BotGroups{},
	}
	*cfg = *global
	return cfg
}

func (c *Config) Load(ctx context.Context, cfg *gcfg.Config) error {
	get, err := cfg.Get(ctx, ".")
	if err != nil {
		return err
	}
	err = get.Scan(c)
	if err != nil {
		return err
	}
	return nil
}
