package config

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// runnerWg 用于等待所有异步 Runner 执行完成
var runnerWg sync.WaitGroup

// Bootstrap 是启动器核心结构，负责统一管理应用生命周期的各个阶段
type Bootstrap struct {
	ctx    context.Context    // 生命周期上下文
	cancel context.CancelFunc // 主动取消上下文的函数
	app    app.IRuntime       // 应用运行时接口，提供模块依赖、全局注入等能力

	before []boot.BootFunc   // 启动前执行的准备逻辑（如日志配置、配置加载）
	boots  []boot.Initialize // 启动阶段初始化的模块（如数据库、缓存、服务注册）
	runs   []boot.BootFunc   // 启动后的主运行逻辑（通常为阻塞服务）
	after  []boot.BootFunc   // 程序结束时的清理逻辑（如释放连接、同步日志等）
}

// NewBootstrap 创建并初始化一个 Bootstrap 启动器实例
func NewBootstrap(ctx context.Context, app app.IRuntime) *Bootstrap {
	c, cancel := context.WithCancel(ctx)
	return &Bootstrap{
		ctx:    c,
		cancel: cancel,
		app:    app,
	}
}

// Before 设置启动前执行的 BootFunc 阶段，通常用于加载配置、日志等
func (b *Bootstrap) Before(fns ...boot.BootFunc) boot.IBootstrap {
	b.before = fns
	return b
}

// Process 设置模块初始化阶段（可用于注册服务、连接数据库等）
func (b *Bootstrap) Process(boots ...boot.Initialize) boot.IBootstrap {
	b.boots = boots
	return b
}

// Runner 设置主运行逻辑（如监听 HTTP 服务、消费者、计划任务等）
func (b *Bootstrap) Runner(fns ...boot.BootFunc) boot.IBootstrap {
	b.runs = fns
	return b
}

// After 设置关闭前的清理逻辑（如关闭连接、日志收尾等）
func (b *Bootstrap) After(fns ...boot.BootFunc) boot.IBootstrap {
	b.after = fns
	return b
}

// Run 启动并管理完整应用生命周期：准备、初始化、运行、等待、清理
func (b *Bootstrap) Run() error {
	// 监听系统中断信号（如 Ctrl+C）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 执行启动准备与初始化
	if err := b.startup(); err != nil {
		return err
	}

	// 阻塞等待退出信号或上下文取消
	select {
	case sig := <-sigChan:
		glog.Infof(b.ctx, "Received signal: %v", sig)
		b.cancel()
	case <-b.ctx.Done():
		glog.Info(b.ctx, "Context canceled")
	}

	glog.Info(b.ctx, "waitForRunners...")
	// 等待所有运行中的 goroutine 正常结束
	b.waitForRunners()

	glog.Info(b.ctx, "cleanup...")
	// 执行清理阶段
	return b.cleanup()
}

// startup 执行 Before、Init、Run 阶段
func (b *Bootstrap) startup() error {
	if err := b.execPhase("before", b.before); err != nil {
		return err
	}
	if err := b.initModules(); err != nil {
		return err
	}
	b.runRunners()
	return nil
}

// initModules 执行所有模块的初始化逻辑（Process）
func (b *Bootstrap) initModules() error {
	for _, mod := range b.boots {
		if err := mod.Init(b.ctx, b.app); err != nil {
			glog.Errorf(b.ctx, "Init module [%s] failed: %v", mod.String(), err)
			return err
		}
	}
	return nil
}

// runRunners 启动所有注册的 Runner 函数（通常为异步阻塞服务）
func (b *Bootstrap) runRunners() {
	for i, fn := range b.runs {
		runnerWg.Add(1)
		go func(i int, f boot.BootFunc) {
			defer runnerWg.Done()

			// 防止单个 goroutine panic 导致进程崩溃
			defer func() {
				if r := recover(); r != nil {
					glog.Errorf(b.ctx, "Runner #%d panic: %v", i, r)
					b.cancel()
				}
			}()

			if err := f(b.ctx, b.app); err != nil {
				glog.Errorf(b.ctx, "Runner #%d error: %v", i, err)
				b.cancel()
			}
		}(i, fn)
	}
}

// waitForRunners 等待所有 Runner 正常退出
func (b *Bootstrap) waitForRunners() {
	runnerWg.Wait()
}

// cleanup 执行 After 阶段的清理逻辑
func (b *Bootstrap) cleanup() error {
	return b.execPhase("after", b.after)
}

// execPhase 用于统一执行指定阶段的 BootFunc 列表，并添加统一日志
func (b *Bootstrap) execPhase(name string, fns []boot.BootFunc) error {
	for i, fn := range fns {
		glog.Infof(b.ctx, "[%s #%d] execPhase...\n", name, i)
		if err := fn(b.ctx, b.app); err != nil {
			glog.Errorf(b.ctx, "[%s #%d] failed: %v", name, i, err)
			return err
		}
	}
	return nil
}
