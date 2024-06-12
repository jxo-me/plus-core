package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
	"os"
	"os/signal"
	"syscall"
)

type Bootstrap struct {
	ctx    context.Context
	cancel      context.CancelFunc
	app    app.IRuntime
	before []boot.BootFunc
	boots  []boot.Initialize
	after  []boot.BootFunc
}

func NewBootstrap(ctx context.Context, app app.IRuntime) *Bootstrap {
	ctx, cancel := context.WithCancel(ctx)
	return &Bootstrap{
		ctx:    ctx,
		cancel: cancel,
		app:    app,
		before: make([]boot.BootFunc, 0),
		boots:  make([]boot.Initialize, 0),
		after:  make([]boot.BootFunc, 0),
	}
}

func (b *Bootstrap) runBootstrap() error {
	var err error
	for _, beforeFunc := range b.before {
		err = beforeFunc(b.ctx, b.app)
		if err != nil {
			glog.Error(b.ctx, fmt.Sprintf("run bootstrap beforeFunc error: %v", err))
			return err
		}
	}
	for i := range b.boots {
		err = b.boots[i].Init(b.ctx, b.app)
		if err != nil {
			glog.Error(b.ctx, fmt.Sprintf("run bootstrap initFunc name: %s error: %v", b.boots[i].String(), err))
			return err
		}
	}
	for _, afterFunc := range b.after {
		err = afterFunc(b.ctx, b.app)
		if err != nil {
			glog.Error(b.ctx, fmt.Sprintf("run bootstrap afterFunc error: %v", err))
			return err
		}
	}

	return err
}

func (b *Bootstrap) Before(before ...boot.BootFunc) boot.IBootstrap {
	b.before = before
	return b
}

func (b *Bootstrap) Process(boots ...boot.Initialize) boot.IBootstrap {
	b.boots = boots
	return b
}

func (b *Bootstrap) After(after ...boot.BootFunc) boot.IBootstrap {
	b.after = after
	return b
}

func (b *Bootstrap) Run() error {
	return b.runBootstrap()
}

// handleSignal sets up signal handling to gracefully shut down the application
func (b *Bootstrap) handleSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("plus app received signal: %v, shutting down...\n", sig)
		b.cancel()
	}()
}
