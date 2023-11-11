package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
)

type BootFunc func(ctx context.Context, app app.IRuntime) error

type Bootstrap struct {
	ctx    context.Context
	app    app.IRuntime
	before []BootFunc
	boots  []boot.Initialize
	after  []BootFunc
}

func NewBoot(ctx context.Context, app app.IRuntime) *Bootstrap {
	return &Bootstrap{
		ctx:    ctx,
		app:    app,
		before: make([]BootFunc, 0),
		boots:  make([]boot.Initialize, 0),
		after:  make([]BootFunc, 0),
	}
}

func (b *Bootstrap) runBootstrap() error {
	var err error
	for _, bootFunc := range b.before {
		err = bootFunc(b.ctx, b.app)
		if err != nil {
			glog.Error(b.ctx, fmt.Sprintf("run bootstrap before error: %v", err))
			return err
		}
	}
	for i := range b.boots {
		err = b.boots[i].Init(b.ctx, b.app)
		if err != nil {
			glog.Error(b.ctx, fmt.Sprintf("run bootstrap %s error: %v", b.boots[i].String(), err))
			return err
		}
	}
	for _, bootFunc := range b.after {
		err = bootFunc(b.ctx, b.app)
		if err != nil {
			glog.Error(b.ctx, fmt.Sprintf("run bootstrap after error: %v", err))
			return err
		}
	}

	return err
}

func (b *Bootstrap) Before(before ...BootFunc) *Bootstrap {
	b.before = before
	return b
}

func (b *Bootstrap) Process(boots ...boot.Initialize) *Bootstrap {
	b.boots = boots
	return b
}

func (b *Bootstrap) After(after ...BootFunc) *Bootstrap {
	b.after = after
	return b
}

func (b *Bootstrap) Run() error {
	return b.runBootstrap()
}
