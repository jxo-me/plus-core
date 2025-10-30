package config

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/jxo-me/plus-core/core/v2/app"
)

type dummyApp struct{ app.IRuntime }

func TestBootstrap_WatcherReload(t *testing.T) {
	b := NewBootstrap(context.Background(), &dummyApp{})
	var reloadCalled int
	b.OnReload(func(ctx context.Context, a app.IRuntime) error {
		reloadCalled++
		return nil
	})
	b.Watchers(func(ctx context.Context, onChange func()) error {
		go func() { time.Sleep(150 * time.Millisecond); onChange() }()
		<-ctx.Done()
		return nil
	})
	b.Runner(func(ctx context.Context, a app.IRuntime) error { <-ctx.Done(); return nil })
	go func() {
		time.Sleep(300 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if reloadCalled < 1 {
		t.Fatal("reload not called by watcher")
	}
}

func TestBootstrap_SIGHUP_Reload(t *testing.T) {
	b := NewBootstrap(context.Background(), &dummyApp{})
	var reloadCalled int
	b.OnReload(func(ctx context.Context, a app.IRuntime) error { reloadCalled++; return nil })
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
	if reloadCalled < 1 {
		t.Fatal("reload not called by SIGHUP")
	}
}

func TestBootstrap_HotPlug_LoadUnload(t *testing.T) {
	b := NewBootstrap(context.Background(), &dummyApp{})
	var loaded, unloaded bool
	b.RegisterLoader("svcA", func(ctx context.Context, a app.IRuntime) error { loaded = true; return nil })
	b.RegisterUnloader("svcA", func(ctx context.Context, a app.IRuntime) error { unloaded = true; return nil })
	b.Runner(func(ctx context.Context, a app.IRuntime) error { <-ctx.Done(); return nil })
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = b.Load("svcA")
		_ = b.Unload("svcA")
		p, _ := os.FindProcess(os.Getpid())
		time.Sleep(100 * time.Millisecond)
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if !loaded || !unloaded {
		t.Fatalf("hot plug not executed: loaded=%v unloaded=%v", loaded, unloaded)
	}
}
