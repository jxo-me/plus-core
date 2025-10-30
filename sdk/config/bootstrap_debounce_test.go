package config

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/jxo-me/plus-core/core/v2/app"
)

type dummyApp2 struct{ app.IRuntime }

func TestReload_Debounce_MergesEvents(t *testing.T) {
	b := NewBootstrap(context.Background(), &dummyApp2{}).EnableAutoReloadCapabilities().EnableReloadDebounce(120 * time.Millisecond)
	var reloadCount int
	b.OnReload(func(ctx context.Context, a app.IRuntime) error { reloadCount++; return nil })
	b.Runner(func(ctx context.Context, a app.IRuntime) error { <-ctx.Done(); return nil })

	// 快速多次触发（小于去抖窗口）
	for i := 0; i < 10; i++ {
		b.handleReload()
	}
	// 再稍后再触发一组事件，也应合并为第二次
	go func() {
		time.Sleep(50 * time.Millisecond)
		for i := 0; i < 5; i++ {
			b.handleReload()
		}
		// 等待合并并退出
		time.Sleep(300 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()

	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	// 预期：两批事件各合并一次，共约2次；考虑调度偶发，至少应为1次，最多不应等于10+5
	if reloadCount < 1 || reloadCount > 4 {
		t.Fatalf("debounce not effective, count=%d", reloadCount)
	}
}
