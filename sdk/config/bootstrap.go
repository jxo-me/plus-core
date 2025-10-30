package config

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
)

// runnerWg 用于等待所有异步 Runner 执行完成
var runnerWg sync.WaitGroup

// Bootstrap 是启动器核心结构，负责统一管理应用生命周期的各个阶段(支持更自动与高级扩展)
type Bootstrap struct {
	ctx    context.Context    // 生命周期上下文
	cancel context.CancelFunc // 主动取消上下文的函数
	app    app.IRuntime       // 应用运行时接口，提供模块依赖、全局注入等能力

	before []boot.BootFunc   // 启动前执行的准备逻辑（如日志配置、配置加载）
	boots  []boot.Initialize // 启动阶段初始化的模块（如数据库、缓存、服务注册）
	runs   []boot.BootFunc   // 启动后的主运行逻辑（通常为阻塞服务）
	after  []boot.BootFunc   // 程序结束时的清理逻辑（如释放连接、同步日志等）

	runnerErrorPolicy string                        // Runner error 策略："exit"/"ignore"/"restart"
	healthChecks      []func(context.Context) error // 健康检查钩子
	sigReload         chan os.Signal                // 支持 SIGHUP 动态 reload

	// 配置变更监听器，触发时进行reload
	watchers []func(ctx context.Context, onChange func()) error
	// reload 钩子：配置或信号触发后的重载回调
	reload []boot.BootFunc
	// 热插拔：动态加载/卸载服务的回调集合
	loaders   map[string]boot.BootFunc
	unloaders map[string]boot.BootFunc

	// 是否自动对各能力Registry中实现Reloadable的对象执行Reload
	autoReloadCaps bool
	// 可选：提供最新配置的提供器，供能力在重载前应用配置
	configProvider func(ctx context.Context) interface{}

	// 重载事件通道与控制
	reloadChan     chan struct{}
	reloadDebounce time.Duration
	reloadMu       sync.Mutex

	// 可选：能力插件注册中心（通过视图接口避免强耦合）
	abilityReg AbilityRegistryView
}

// 能力注册中心简化视图接口
// 仅用于获取与遍历插件对象

type AbilityRegistryView interface {
	Get(name string) (interface{}, bool)
	List() []string
}

// 本地断言接口，避免跨模块强依赖
// 任何实现 Reload(context.Context) error 的对象都将被自动重载

type capabilityReloadable interface {
	Reload(ctx context.Context) error
}

// 本地断言接口：实现 ApplyConfig(any) error 的对象会先应用配置

type capabilityConfigurable interface {
	ApplyConfig(cfg interface{}) error
}

// NewBootstrap 创建并初始化一个 Bootstrap 启动器实例
func NewBootstrap(ctx context.Context, app app.IRuntime) *Bootstrap {
	c, cancel := context.WithCancel(ctx)
	return &Bootstrap{
		ctx:        c,
		cancel:     cancel,
		app:        app,
		sigReload:  make(chan os.Signal, 1),
		loaders:    make(map[string]boot.BootFunc),
		unloaders:  make(map[string]boot.BootFunc),
		reloadChan: make(chan struct{}, 16),
	}
}

// 可选增强配置
func (b *Bootstrap) WithRunnerErrorPolicy(policy string) *Bootstrap {
	b.runnerErrorPolicy = policy
	return b
}

func (b *Bootstrap) WithHealthCheck(checks ...func(context.Context) error) *Bootstrap {
	b.healthChecks = append(b.healthChecks, checks...)
	return b
}

// WithReloadBuffer 设置reload事件通道缓冲
func (b *Bootstrap) WithReloadBuffer(n int) *Bootstrap {
	if n <= 0 {
		n = 1
	}
	b.reloadChan = make(chan struct{}, n)
	return b
}

// EnableReloadDebounce 启用重载去抖（将短时间内多次变更合并为一次）
func (b *Bootstrap) EnableReloadDebounce(d time.Duration) *Bootstrap {
	b.reloadDebounce = d
	return b
}

// Watchers 注册配置变更监听器。监听器内部检测到变更后应调用 onChange()。
func (b *Bootstrap) Watchers(fns ...func(ctx context.Context, onChange func()) error) *Bootstrap {
	b.watchers = append(b.watchers, fns...)
	return b
}

// OnReload 注册重载钩子，在收到 SIGHUP 或配置变更时被执行。
func (b *Bootstrap) OnReload(fns ...boot.BootFunc) *Bootstrap {
	b.reload = append(b.reload, fns...)
	return b
}

// EnableAutoReloadCapabilities 打开自动对实现 app.Reloadable 的能力对象进行重载
func (b *Bootstrap) EnableAutoReloadCapabilities() *Bootstrap {
	b.autoReloadCaps = true
	return b
}

// WithConfigProvider 设置配置提供器；在自动能力重载时会为实现了 ApplyConfig 的能力先下发配置
func (b *Bootstrap) WithConfigProvider(provider func(ctx context.Context) interface{}) *Bootstrap {
	b.configProvider = provider
	return b
}

// WithAbilityRegistry 注入能力插件注册中心，便于在重载与热插拔时统一编排
func (b *Bootstrap) WithAbilityRegistry(ar AbilityRegistryView) *Bootstrap {
	b.abilityReg = ar
	return b
}

// RegisterLoader 注册热加载回调（按名称）。调用 Load(name) 时执行。
func (b *Bootstrap) RegisterLoader(name string, fn boot.BootFunc) *Bootstrap {
	b.loaders[name] = fn
	return b
}

// RegisterUnloader 注册热卸载回调（按名称）。调用 Unload(name) 时执行。
func (b *Bootstrap) RegisterUnloader(name string, fn boot.BootFunc) *Bootstrap {
	b.unloaders[name] = fn
	return b
}

// Load 触发指定名称的服务加载（热插拔装载）。
func (b *Bootstrap) Load(name string) error {
	if b.abilityReg != nil {
		if p, ok := b.abilityReg.Get(name); ok {
			// 插件可选实现启动接口
			if s, ok := any(p).(interface {
				Start(ctx context.Context) error
			}); ok {
				return s.Start(b.ctx)
			}
		}
	}
	if fn, ok := b.loaders[name]; ok {
		return fn(b.ctx, b.app)
	}
	return nil
}

// Unload 触发指定名称的服务卸载（热插拔卸载）。
func (b *Bootstrap) Unload(name string) error {
	if b.abilityReg != nil {
		if p, ok := b.abilityReg.Get(name); ok {
			// 插件可选实现停止接口
			if s, ok := any(p).(interface {
				Stop(ctx context.Context) error
			}); ok {
				return s.Stop(b.ctx)
			}
		}
	}
	if fn, ok := b.unloaders[name]; ok {
		return fn(b.ctx, b.app)
	}
	return nil
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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer signal.Stop(sigChan)

	if err := b.startup(); err != nil {
		return err
	}
	// startup后阶段健康检查
	for _, f := range b.healthChecks {
		if err := f(b.ctx); err != nil {
			return errors.New("healthcheck failed: " + err.Error())
		}
	}
	done := false
	for !done {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGHUP:
				go b.handleReload()
				continue // reload不退出
			case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
				glog.Infof(b.ctx, "Received exit signal: %v", sig)
				b.cancel()
				done = true
			}
		case <-b.ctx.Done():
			glog.Info(b.ctx, "Context canceled")
			done = true
		}
	}

	b.waitForRunners()
	return b.cleanup()
}

// handleReload 作为唯一编排入口：仅推送事件到重载通道，由worker去抖与串行执行
func (b *Bootstrap) handleReload() {
	select {
	case b.reloadChan <- struct{}{}:
	default:
		// 通道已满，丢弃多余事件，等待worker合并
	}
}

// processReload 串行执行实际重载：ApplyConfig -> Reload -> OnReload
func (b *Bootstrap) processReload() {
	b.reloadMu.Lock()
	defer b.reloadMu.Unlock()
	glog.Info(b.ctx, "processing reload...")
	if b.autoReloadCaps {
		b.reloadCapabilities()
		b.reloadPlugins()
	}
	_ = b.execPhase("reload", b.reload)
}

// reloaderWorker 消费重载事件，带去抖合并
func (b *Bootstrap) reloaderWorker() {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorf(b.ctx, "reloaderWorker panic: %v", r)
		}
	}()
	for {
		select {
		case <-b.ctx.Done():
			return
		case <-b.reloadChan:
			// 去抖窗口
			if b.reloadDebounce > 0 {
				t := time.NewTimer(b.reloadDebounce)
				debouncing := true
				for debouncing {
					select {
					case <-b.ctx.Done():
						if !t.Stop() {
							<-t.C
						}
						return
					case <-b.reloadChan:
						// 吸收更多事件直至定时器到期
					case <-t.C:
						debouncing = false
					}
				}
			}
			b.processReload()
		}
	}
}

// reloadCapabilities 遍历已注册的能力对象，对实现 app.Reloadable 的执行 Reload
func (b *Bootstrap) reloadCapabilities() {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorf(b.ctx, "reloadCapabilities panic: %v", r)
		}
	}()
	apply := func(v interface{}) {
		if b.configProvider != nil {
			if c, ok := v.(capabilityConfigurable); ok {
				_ = c.ApplyConfig(b.configProvider(b.ctx))
			}
		}
		if r, ok := v.(capabilityReloadable); ok {
			_ = r.Reload(b.ctx)
		}
	}
	// cache
	for _, v := range b.app.CacheRegistry().GetAll() {
		apply(v)
	}
	// queue
	for _, v := range b.app.QueueRegistry().GetAll() {
		apply(v)
	}
	// task service
	for _, v := range b.app.TaskRegister().GetAll() {
		apply(v)
	}
	// memory task service
	for _, v := range b.app.MemoryTaskRegister().GetAll() {
		apply(v)
	}
	// jwt
	for _, v := range b.app.JwtRegister().GetAll() {
		apply(v)
	}
	// verifier
	for _, v := range b.app.VerifyRegistry().GetAll() {
		apply(v)
	}
}

// reloadPlugins 针对AbilityRegistry内的插件执行 ApplyConfig 与 Reload
func (b *Bootstrap) reloadPlugins() {
	if b.abilityReg == nil {
		return
	}
	apply := func(v interface{}) {
		if b.configProvider != nil {
			if c, ok := v.(capabilityConfigurable); ok {
				_ = c.ApplyConfig(b.configProvider(b.ctx))
			}
		}
		if r, ok := v.(capabilityReloadable); ok {
			_ = r.Reload(b.ctx)
		}
	}
	for _, name := range b.abilityReg.List() {
		if v, ok := b.abilityReg.Get(name); ok {
			apply(v)
		}
	}
}

// startup 执行 Before、Init、Run 阶段
func (b *Bootstrap) startup() error {
	if err := b.execPhase("before", b.before); err != nil {
		return err
	}
	if err := b.initModules(); err != nil {
		return err
	}
	// 启动重载worker
	runnerWg.Add(1)
	go func() {
		defer runnerWg.Done()
		b.reloaderWorker()
	}()
	b.runRunners()
	b.runWatchers()
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
			defer func() {
				if r := recover(); r != nil {
					glog.Errorf(b.ctx, "Runner #%d panic: %v", i, r)
					if b.runnerErrorPolicy == "restart" {
						glog.Warningf(b.ctx, "Restarting runner #%d ...", i)
						go b.restartRunner(i, f, 1)
						return
					}
					if b.runnerErrorPolicy != "ignore" {
						b.cancel()
					}
				}
			}()
			if err := f(b.ctx, b.app); err != nil {
				glog.Errorf(b.ctx, "Runner #%d error: %v", i, err)
				if b.runnerErrorPolicy == "restart" {
					glog.Warningf(b.ctx, "Restarting runner #%d ...", i)
					go b.restartRunner(i, f, 1)
					return
				}
				if b.runnerErrorPolicy != "ignore" {
					b.cancel()
				}
			}
		}(i, fn)
	}
}

// restartRunner 支持Runner崩溃时自动重启，带指数退避
func (b *Bootstrap) restartRunner(i int, f boot.BootFunc, times int) {
	glog.Warningf(b.ctx, "Restarting #%d attempt %d", i, times)
	d := time.Duration(times*times) * time.Second
	time.Sleep(d)
	runnerWg.Add(1)
	go func() {
		defer runnerWg.Done()
		if err := f(b.ctx, b.app); err != nil && times < 5 {
			b.restartRunner(i, f, times+1)
		}
	}()
}

// waitForRunners 等待所有 Runner 正常退出
func (b *Bootstrap) waitForRunners() {
	glog.Info(b.ctx, "waitForRunners...")
	runnerWg.Wait()
}

// cleanup 执行 After 阶段的清理逻辑
func (b *Bootstrap) cleanup() error {
	return b.execPhase("after", b.after)
}

// execPhase 用于统一执行指定阶段的 BootFunc 列表，并添加统一日志
func (b *Bootstrap) execPhase(name string, fns []boot.BootFunc) error {
	glog.Info(b.ctx, "cleanup...")
	for i, fn := range fns {
		glog.Infof(b.ctx, "[%s #%d] execPhase...\n", name, i)
		if err := fn(b.ctx, b.app); err != nil {
			glog.Errorf(b.ctx, "[%s #%d] failed: %v", name, i, err)
			return err
		}
	}
	return nil
}

// 启动配置变更监听器
func (b *Bootstrap) runWatchers() {
	for i, w := range b.watchers {
		runnerWg.Add(1)
		go func(i int, wf func(ctx context.Context, onChange func()) error) {
			defer runnerWg.Done()
			onChange := func() {
				glog.Infof(b.ctx, "config changed by watcher #%d, triggering reload", i)
				b.handleReload()
			}
			// 忽略 watcher 错误，仅记录
			_ = wf(b.ctx, onChange)
		}(i, w)
	}
}
