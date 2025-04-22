package config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var runnerWg sync.WaitGroup

type Bootstrap struct {
	ctx    context.Context
	cancel context.CancelFunc
	app    app.IRuntime
	before []boot.BootFunc
	boots  []boot.Initialize
	runs   []boot.BootFunc
	after  []boot.BootFunc
}

func NewBootstrap(ctx context.Context, app app.IRuntime) *Bootstrap {
	c, cancel := context.WithCancel(ctx)
	return &Bootstrap{
		ctx:    c,
		cancel: cancel,
		app:    app,
	}
}

func (b *Bootstrap) Before(before ...boot.BootFunc) boot.IBootstrap {
	b.before = before
	return b
}

func (b *Bootstrap) Process(boots ...boot.Initialize) boot.IBootstrap {
	b.boots = boots
	return b
}

func (b *Bootstrap) Runner(runs ...boot.BootFunc) boot.IBootstrap {
	b.runs = runs
	return b
}

func (b *Bootstrap) After(after ...boot.BootFunc) boot.IBootstrap {
	b.after = after
	return b
}

func (b *Bootstrap) Run() error {
	// 捕捉系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动逻辑
	if err := b.runStartup(); err != nil {
		return err
	}

	// 等待退出信号
	select {
	case sig := <-sigChan:
		glog.Infof(b.ctx, "Received signal: %v, shutting down...", sig)
		b.cancel()
	case <-b.ctx.Done():
		glog.Info(b.ctx, "Context canceled, exiting...")
	}

	// 等待 Runner() 的 goroutine 完成
	b.waitForRunners()

	// 执行 After() 收尾逻辑
	return b.runAfter()
}

func (b *Bootstrap) runStartup() error {
	// 1. Before
	for _, fn := range b.before {
		if err := fn(b.ctx, b.app); err != nil {
			glog.Error(b.ctx, fmt.Sprintf("beforeFunc error: %v", err))
			return err
		}
	}

	// 2. Init
	for _, boot := range b.boots {
		if err := boot.Init(b.ctx, b.app); err != nil {
			glog.Error(b.ctx, fmt.Sprintf("initFunc [%s] error: %v", boot.String(), err))
			return err
		}
	}

	// 3. Run
	b.runRunners()

	return nil
}

func (b *Bootstrap) runRunners() {
	for _, fn := range b.runs {
		runnerWg.Add(1)
		go func(f boot.BootFunc) {
			defer runnerWg.Done()
			if err := f(b.ctx, b.app); err != nil {
				glog.Warning(b.ctx, fmt.Sprintf("runFunc error: %v", err))
				b.cancel()
			}
		}(fn)
	}
}

func (b *Bootstrap) waitForRunners() {
	runnerWg.Wait()
}

func (b *Bootstrap) runAfter() error {
	for _, fn := range b.after {
		if err := fn(context.Background(), b.app); err != nil {
			glog.Error(b.ctx, fmt.Sprintf("afterFunc error: %v", err))
			return err
		}
	}
	return nil
}
