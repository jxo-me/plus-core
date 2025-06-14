package logger

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func Infof(ctx context.Context, format string, v ...interface{}) {
	g.Log().Infof(ctx, format, v...)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	g.Log().Warningf(ctx, format, v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	g.Log().Errorf(ctx, format, v...)
}
