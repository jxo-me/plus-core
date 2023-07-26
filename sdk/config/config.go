package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/pkg/tus"
	"github.com/jxo-me/plus-core/sdk/pkg/ws"
)

var (
	ExtendConfig interface{}
	insSetting   = Settings{}
)

// Settings 兼容原先的配置结构
type Settings struct {
	srv       *ghttp.Server
	cfg       *gcfg.Config
	config    Config `yaml:"config"`
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

// Config 配置集合
type Config struct {
	Jwt     map[string]*Jwt `yaml:"jwt"`
	Cache   *Cache          `yaml:"cache"`
	Queue   *Queue          `yaml:"queue"`
	Locker  *Locker         `yaml:"locker"`
	Extend  interface{}     `yaml:"extend"`
	Tus     tus.Config      `yaml:"tus"`
	Ws      *ws.Config      `yaml:"ws"`
	Metrics *Metrics        `yaml:"metrics"`
}

// Bootstrap 载入启动配置文件
func (e *Settings) Bootstrap(ctx context.Context, fs ...Initialize) {
	e.config = Config{
		Jwt:    map[string]*Jwt{},
		Cache:  CacheConfig(),
		Queue:  QueueConfig(),
		Extend: ExtendConfig,
		Locker: LockerConfig(),
		Ws:     &ws.Config{},
	}
	e.callbacks = fs
	e.runCallback(ctx)
}

func (e *Settings) SetCfg(cf *gcfg.Config) *Settings {
	e.cfg = cf
	return e
}

func (e *Settings) Cfg() *gcfg.Config {
	return e.cfg
}

func (e *Settings) SetConfig(c Config) *Settings {
	e.config = c
	return e
}

func (e *Settings) Config() *Config {
	return &e.config
}

func (e *Settings) SetSrv(srv *ghttp.Server) *Settings {
	e.srv = srv
	return e
}

func (e *Settings) Srv() *ghttp.Server {
	return e.srv
}

func (e *Settings) SetJwt(module string, jwt *Jwt) *Settings {
	e.Config().Jwt[module] = jwt
	return e
}

func (e *Settings) GetJwt(module string) *Jwt {
	if j, ok := e.Config().Jwt[module]; ok {
		return j
	}
	return nil
}

func (e *Settings) GetTus() tus.Config {
	return e.Config().Tus
}
