package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	ExtendConfig interface{}
	insSetting   = Settings{}
)

type Initialize interface {
	String() string
	Init(ctx context.Context, s *Settings) error
}

// Settings 兼容原先的配置结构
type Settings struct {
	srv       *ghttp.Server
	cfg       *gcfg.Config
	Config    Config `yaml:"config"`
	callbacks []Initialize
}

func Setting() *Settings {
	return &insSetting
}

func (e *Settings) runCallback(ctx context.Context) {
	for i := range e.callbacks {
		err := e.callbacks[i].Init(ctx, e)
		if err != nil {
			glog.Error(ctx, fmt.Sprintf("runCallback %s error: %v", e.callbacks[i].String(), err))
		}
	}
}

func (e *Settings) Init(ctx context.Context) {
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

// Bootstrap 载入启动配置文件
func (e *Settings) Bootstrap(ctx context.Context, fs ...Initialize) {
	e.Config = Config{
		Jwt:    JwtConfig,
		Cache:  CacheConfig,
		Queue:  QueueConfig,
		Extend: ExtendConfig,
		Locker: LockerConfig,
	}
	e.callbacks = fs
	e.Init(ctx)
}

func (e *Settings) SetCfg(cf *gcfg.Config) *Settings {
	e.cfg = cf
	return e
}

func (e *Settings) Cfg() *gcfg.Config {
	return e.cfg
}

func (e *Settings) SetSrv(srv *ghttp.Server) *Settings {
	e.srv = srv
	return e
}

func (e *Settings) Srv() *ghttp.Server {
	return e.srv
}
