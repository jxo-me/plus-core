package boot

import (
	"context"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/queue"
)

type BootFunc func(ctx context.Context, app app.IRuntime) error

type Initialize interface {
	String() string
	Init(ctx context.Context, app app.IRuntime) error
}

type QueueInitialize interface {
	Initialize
	GetQueue(ctx context.Context) (map[string]queue.IQueue, error)
}

type IBootstrap interface {
	Before(before ...BootFunc) IBootstrap
	Process(boots ...Initialize) IBootstrap
	After(after ...BootFunc) IBootstrap
	Runner(runs ...BootFunc) IBootstrap
	Run() error
}
