package config

import (
	"context"
	"testing"
	"time"
)

type arView struct{ m map[string]interface{} }

func (a *arView) Get(name string) (interface{}, bool) { v, ok := a.m[name]; return v, ok }
func (a *arView) List() []string {
	ks := make([]string, 0, len(a.m))
	for k := range a.m {
		ks = append(ks, k)
	}
	return ks
}

type pluginX struct {
	applied  int
	reloaded int
	started  int
	stopped  int
}

func (p *pluginX) ApplyConfig(cfg interface{}) error {
	if m, ok := cfg.(map[string]any); ok {
		if v, ok := m["v"].(int); ok {
			p.applied = v
		}
	}
	return nil
}
func (p *pluginX) Reload(ctx context.Context) error { p.reloaded++; return nil }
func (p *pluginX) Start(ctx context.Context) error  { p.started++; return nil }
func (p *pluginX) Stop(ctx context.Context) error   { p.stopped++; return nil }

func TestAbilityRegistry_Integration(t *testing.T) {
	pl := &pluginX{}
	reg := &arView{m: map[string]interface{}{"svc": pl}}
	b := NewBootstrap(context.Background(), nil).WithAbilityRegistry(reg).EnableAutoReloadCapabilities().WithConfigProvider(func(ctx context.Context) interface{} { return map[string]any{"v": 7} })
	// 使用内置loader/unloader优先走AbilityRegistry的Start/Stop
	if err := b.Load("svc"); err != nil {
		t.Fatalf("load err: %v", err)
	}
	if err := b.Unload("svc"); err != nil {
		t.Fatalf("unload err: %v", err)
	}
	// 触发一次重载
	b.handleReload()
	// 等待worker处理
	time.Sleep(50 * time.Millisecond)
	if pl.applied != 7 || pl.reloaded < 1 || pl.started < 1 || pl.stopped < 1 {
		t.Fatalf("plugin not integrated: applied=%d reloaded=%d started=%d stopped=%d", pl.applied, pl.reloaded, pl.started, pl.stopped)
	}
}
