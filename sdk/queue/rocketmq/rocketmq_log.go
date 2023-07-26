package rocketmq

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"strings"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type QLoger struct {
	Ctx    context.Context
	Logger *glog.Logger
}

func (q *QLoger) OutputPath(path string) (err error) {
	if "" == path {
		return
	}
	q.Logger.Path(path)
	return nil
}

func (q *QLoger) Debug(msg string, fields map[string]interface{}) {
	q.Logger.Debug(q.Ctx, msg, fields)
}

func (q *QLoger) Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	q.Logger.Info(q.Ctx, msg, fields)
}

func (q *QLoger) Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	q.Logger.Warning(q.Ctx, msg, fields)
}

func (q *QLoger) Error(msg string, fields map[string]interface{}) {
	q.Logger.Error(q.Ctx, msg, fields)
}

func (q *QLoger) Fatal(msg string, fields map[string]interface{}) {
	q.Logger.Fatal(q.Ctx, msg, fields)
}

func (q *QLoger) Level(level string) {
	switch strings.ToLower(level) {
	case "debug":
		q.Logger.SetLevel(DebugLevel)
	case "warn":
		q.Logger.SetLevel(WarnLevel)
	case "error":
		q.Logger.SetLevel(ErrorLevel)
	case "fatal":
		q.Logger.SetLevel(FatalLevel)
	default:
		q.Logger.SetLevel(InfoLevel)
	}
}
