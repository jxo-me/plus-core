package config

import (
	"context"
	"github.com/gogf/gf/v2/os/gcfg"
	"log"
)

var (
	ExtendConfig interface{}
	_cfg         *Settings
)

// Settings 兼容原先的配置结构
type Settings struct {
	Config    *gcfg.Config
	Settings  Config `yaml:"settings"`
	Cache     *Cache `yaml:"cache"`
	Queue     *Queue `yaml:"queue"`
	callbacks []func(ctx context.Context, cf *gcfg.Config)
}

func (e *Settings) runCallback(ctx context.Context) {
	for i := range e.callbacks {
		e.callbacks[i](ctx, e.Config)
	}
}

func (e *Settings) Init(ctx context.Context) {
	e.init(ctx)
	log.Println("!!! config init")
}

func (e *Settings) init(ctx context.Context) {
	e.runCallback(ctx)
}

// Config 配置集合
type Config struct {
	Jwt    *Jwt        `yaml:"jwt"`
	Cache  *Cache      `yaml:"cache"`
	Queue  *Queue      `yaml:"queue"`
	Locker *Locker     `yaml:"locker"`
	Extend interface{} `yaml:"extend"`
}

// Setup 载入配置文件
func Setup(ctx context.Context, s *gcfg.Config, fs ...func(ctx context.Context, s *gcfg.Config)) {
	_cfg = &Settings{
		Config: s,
		Settings: Config{
			Jwt:    JwtConfig,
			Cache:  CacheConfig,
			Queue:  QueueConfig,
			Extend: ExtendConfig,
			Locker: LockerConfig,
		},
		callbacks: fs,
	}

	_cfg.Init(ctx)
}
