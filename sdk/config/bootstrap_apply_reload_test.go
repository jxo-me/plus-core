package config

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/jxo-me/plus-core/core/v2/queue"
	reg "github.com/jxo-me/plus-core/core/v2/registry"
	"github.com/jxo-me/plus-core/core/v2/task"
	"github.com/jxo-me/plus-core/pkg/v2/security/auth"
)

type arRuntime struct {
	app.IRuntime
	c  reg.IRegistry[cache.ICache]
	q  reg.IRegistry[queue.IQueue]
	t  reg.IRegistry[task.TasksService]
	mt reg.IRegistry[task.MemoryService]
	j  reg.IRegistry[*jwt.GfJWTMiddleware]
	v  reg.IRegistry[*auth.Verifier]
}

func (a *arRuntime) CacheRegistry() reg.IRegistry[cache.ICache]            { return a.c }
func (a *arRuntime) QueueRegistry() reg.IRegistry[queue.IQueue]            { return a.q }
func (a *arRuntime) TaskRegister() reg.IRegistry[task.TasksService]        { return a.t }
func (a *arRuntime) MemoryTaskRegister() reg.IRegistry[task.MemoryService] { return a.mt }
func (a *arRuntime) JwtRegister() reg.IRegistry[*jwt.GfJWTMiddleware]      { return a.j }
func (a *arRuntime) VerifyRegistry() reg.IRegistry[*auth.Verifier]         { return a.v }

// 简单内存Registry实现

type mr[T any] struct{ m map[string]T }

func (m *mr[T]) Register(name string, v T) error {
	if m.m == nil {
		m.m = make(map[string]T)
	}
	m.m[name] = v
	return nil
}
func (m *mr[T]) Unregister(name string)        { delete(m.m, name) }
func (m *mr[T]) IsRegistered(name string) bool { _, ok := m.m[name]; return ok }
func (m *mr[T]) Get(name string) T             { return m.m[name] }
func (m *mr[T]) GetAll() map[string]T          { return m.m }

// 能力对象：实现ApplyConfig与Reload，并实现 cache.ICache 接口以便注册进 CacheRegistry

type capCache struct {
	applied  interface{}
	reloaded int
}

func (c *capCache) ApplyConfig(cfg interface{}) error { c.applied = cfg; return nil }
func (c *capCache) Reload(ctx context.Context) error  { c.reloaded++; return nil }

// 实现 cache.ICache 所有方法（空实现用于测试）
func (c *capCache) String() string                                         { return "capCache" }
func (c *capCache) Get(ctx context.Context, key string) (*gvar.Var, error) { return nil, nil }
func (c *capCache) Set(ctx context.Context, key string, val interface{}, expire int) error {
	return nil
}
func (c *capCache) Del(ctx context.Context, key string) error                      { return nil }
func (c *capCache) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) { return nil, nil }
func (c *capCache) HashDel(ctx context.Context, hk, key string) error              { return nil }
func (c *capCache) HashSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error) {
	return 0, nil
}
func (c *capCache) HashVals(ctx context.Context, key string) (gvar.Vars, error)   { return nil, nil }
func (c *capCache) HashLen(ctx context.Context, key string) (int64, error)        { return 0, nil }
func (c *capCache) HashGetAll(ctx context.Context, key string) (*gvar.Var, error) { return nil, nil }
func (c *capCache) HashMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	return nil
}
func (c *capCache) HashMGet(ctx context.Context, key string, fields ...string) (gvar.Vars, error) {
	return nil, nil
}
func (c *capCache) ListLPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, nil
}
func (c *capCache) ListRange(ctx context.Context, key string, start, stop int64) (gvar.Vars, error) {
	return nil, nil
}
func (c *capCache) ListRPop(ctx context.Context, key string, count ...int) (*gvar.Var, error) {
	return nil, nil
}
func (c *capCache) Increase(ctx context.Context, key string) (int64, error)         { return 0, nil }
func (c *capCache) Decrease(ctx context.Context, key string) (int64, error)         { return 0, nil }
func (c *capCache) Expire(ctx context.Context, key string, dur time.Duration) error { return nil }

func TestApplyConfigThenReload_ByWatcher(t *testing.T) {
	c := &capCache{}
	cReg := &mr[cache.ICache]{m: map[string]cache.ICache{"x": c}}
	ar := &arRuntime{c: cReg, q: &mr[queue.IQueue]{}, t: &mr[task.TasksService]{}, mt: &mr[task.MemoryService]{}, j: &mr[*jwt.GfJWTMiddleware]{}, v: &mr[*auth.Verifier]{}}
	b := NewBootstrap(context.Background(), ar).EnableAutoReloadCapabilities().WithConfigProvider(func(ctx context.Context) interface{} { return map[string]any{"k": 1} })
	b.Watchers(func(ctx context.Context, onChange func()) error {
		go func() { time.Sleep(120 * time.Millisecond); onChange() }()
		<-ctx.Done()
		return nil
	})
	b.Runner(func(ctx context.Context, a app.IRuntime) error { <-ctx.Done(); return nil })
	go func() {
		time.Sleep(240 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if c.applied == nil {
		t.Fatal("ApplyConfig not called")
	}
	if c.reloaded < 1 {
		t.Fatal("Reload not called")
	}
}

func TestApplyConfigThenReload_BySIGHUP(t *testing.T) {
	c := &capCache{}
	cReg := &mr[cache.ICache]{m: map[string]cache.ICache{"x": c}}
	ar := &arRuntime{c: cReg, q: &mr[queue.IQueue]{}, t: &mr[task.TasksService]{}, mt: &mr[task.MemoryService]{}, j: &mr[*jwt.GfJWTMiddleware]{}, v: &mr[*auth.Verifier]{}}
	b := NewBootstrap(context.Background(), ar).EnableAutoReloadCapabilities().WithConfigProvider(func(ctx context.Context) interface{} { return map[string]any{"k": 2} })
	b.Runner(func(ctx context.Context, a app.IRuntime) error { <-ctx.Done(); return nil })
	go func() {
		time.Sleep(120 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGHUP)
		time.Sleep(120 * time.Millisecond)
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	m, ok := c.applied.(map[string]any)
	if !ok || m["k"].(int) != 2 {
		t.Fatal("ApplyConfig value incorrect")
	}
	if c.reloaded < 1 {
		t.Fatal("Reload not called")
	}
}
