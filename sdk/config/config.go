package config

import (
	"github.com/gogf/gf/v2/os/gcfg"
	"log"
)

var (
	ExtendConfig interface{}
	_cfg         *Settings
)

// Settings 兼容原先的配置结构
type Settings struct {
	Settings  Config `yaml:"settings"`
	Cache     *Cache `yaml:"cache"`
	Queue     *Queue `yaml:"queue"`
	callbacks []func()
}

func (e *Settings) runCallback() {
	for i := range e.callbacks {
		e.callbacks[i]()
	}
}

func (e *Settings) Init() {
	e.init()
	log.Println("!!! config init")
}

func (e *Settings) init() {
	e.runCallback()
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
func Setup(s *gcfg.Config, fs ...func()) {
	_cfg = &Settings{
		Settings: Config{
			Jwt:    JwtConfig,
			Cache:  CacheConfig,
			Queue:  QueueConfig,
			Extend: ExtendConfig,
			Locker: LockerConfig,
		},
		callbacks: fs,
	}

	_cfg.Init()
}
