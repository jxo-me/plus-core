package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/boot"
	"github.com/jxo-me/plus-core/sdk/v2"
)

func runCallback(ctx context.Context, callbacks []boot.Initialize) {
	for i := range callbacks {
		err := callbacks[i].Init(ctx, sdk.Runtime)
		if err != nil {
			glog.Error(ctx, fmt.Sprintf("runCallback %s error: %v", callbacks[i].String(), err))
		}
	}
}

// Bootstrap 载入启动配置文件
func Bootstrap(ctx context.Context, fs ...boot.Initialize) {
	runCallback(ctx, fs)
}
