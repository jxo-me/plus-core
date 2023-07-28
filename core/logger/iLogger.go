package logger

import (
	"context"
)

// ILogger is the interface to send logs to. It can be set using
type ILogger interface {
	Fatalf(context.Context, string, ...interface{})
	Errorf(context.Context, string, ...interface{})
	Warningf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Debugf(context.Context, string, ...interface{})
	Noticef(context.Context, string, ...interface{})
}
