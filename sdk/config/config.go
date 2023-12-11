package config

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/jxo-me/plus-core/pkg/v2/erlang"
	"github.com/jxo-me/plus-core/sdk/v2/send/telegram"
	"sync"
)

var (
	global = &Config{
		Database: make(map[string]*gdb.ConfigNode),
		Redis:    make(map[string]*RedisOptions),
		Server:   &ghttp.ServerConfig{},
		Settings: &SettingOptions{
			Queue: &QueueGroups{
				Rocketmq: make(map[string]*RocketmqOptions),
				Rabbitmq: make(map[string]*RabbitmqOptions),
			},
			Auth:        make(map[string]*JwtAuth),
			FailedLimit: make(map[string]*FailedLimitOptions),
			Uploads:     &UploadGroups{Tus: make(map[string]*TusOptions)},
			RateLimit:   &RateLimitOptions{},
		},
		Bot:      &BotGroups{},
		Telegram: &telegram.SendConf{},
		Erlang:   make(map[string]*erlang.NodeConfig),
	}
	globalMux sync.RWMutex
)

func Global() *Config {
	globalMux.RLock()
	defer globalMux.RUnlock()

	cfg := &Config{
		Database: make(map[string]*gdb.ConfigNode),
		Redis:    make(map[string]*RedisOptions),
		Server:   &ghttp.ServerConfig{},
		Settings: &SettingOptions{
			Queue: &QueueGroups{
				Rocketmq: make(map[string]*RocketmqOptions),
				Rabbitmq: make(map[string]*RabbitmqOptions),
			},
			Auth:        make(map[string]*JwtAuth),
			FailedLimit: make(map[string]*FailedLimitOptions),
			Uploads:     &UploadGroups{Tus: make(map[string]*TusOptions)},
			RateLimit:   &RateLimitOptions{},
		},
		Bot:      &BotGroups{},
		Telegram: &telegram.SendConf{},
		Erlang:   make(map[string]*erlang.NodeConfig),
	}
	*cfg = *global
	return cfg
}

func (c *Config) Load(ctx context.Context, cfg *gcfg.Config) error {
	get, err := cfg.Get(ctx, ".")
	if err != nil {
		return err
	}
	err = get.Scan(global)
	if err != nil {
		return err
	}
	return nil
}
