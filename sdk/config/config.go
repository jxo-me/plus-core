package config

import (
	"context"
	"github.com/gogf/gf/v2/os/gcfg"
	"sync"
)

var (
	global    = &Config{}
	globalMux sync.RWMutex
)

func Global() *Config {
	globalMux.RLock()
	defer globalMux.RUnlock()

	cfg := &Config{}
	*cfg = *global
	return cfg
}

func (c *Config) Load(ctx context.Context, cfg *gcfg.Config) error {
	get, err := cfg.Get(ctx, "settings")
	if err != nil {
		return err
	}
	err = get.Scan(&c.Settings)
	if err != nil {
		return err
	}
	get, err = cfg.Get(ctx, "database")
	if err != nil {
		return err
	}
	err = get.Scan(&c.Database)
	if err != nil {
		return err
	}
	get, err = cfg.Get(ctx, "redis")
	if err != nil {
		return err
	}
	err = get.Scan(&c.Redis)
	if err != nil {
		return err
	}
	get, err = cfg.Get(ctx, "server")
	if err != nil {
		return err
	}
	err = get.Scan(&c.Server)
	if err != nil {
		return err
	}
	get, err = cfg.Get(ctx, "bot")
	if err != nil {
		return err
	}
	err = get.Scan(&c.Bot)
	if err != nil {
		return err
	}
	return nil
}
