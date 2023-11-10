package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/boot"
	"github.com/jxo-me/plus-core/sdk/v2"
)

var (
	insSetting = Settings{}
)

// Settings 兼容原先的配置结构
type Settings struct {
	callbacks []boot.Initialize
}

func Setting() *Settings {
	return &insSetting
}

func (e *Settings) runCallback(ctx context.Context) {
	for i := range e.callbacks {
		err := e.callbacks[i].Init(ctx, sdk.Runtime.ConfigRegister().Get(DefaultGroupName))
		if err != nil {
			glog.Error(ctx, fmt.Sprintf("runCallback %s error: %v", e.callbacks[i].String(), err))
		}
	}
}

// Config 配置集合
type Config struct {
	Cache  *Cache      `yaml:"cache"`
	Extend interface{} `yaml:"extend"`
}

// Bootstrap 载入启动配置文件
func (e *Settings) Bootstrap(ctx context.Context, fs ...boot.Initialize) {
	e.callbacks = fs
	e.runCallback(ctx)
}
