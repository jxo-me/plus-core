package config

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
)

// mockApp 用于模拟 IRuntime
var _ app.IRuntime = (*mockApp)(nil)

type mockApp struct{}

func (m *mockApp) ArpcServerRegistry() interface{}            { return nil }
func (m *mockApp) ArpcClientRegistry() interface{}            { return nil }
func (m *mockApp) BotRegistry() interface{}                   { return nil }
func (m *mockApp) CacheRegistry() interface{}                 { return nil }
func (m *mockApp) CasBinRegistry() interface{}                { return nil }
func (m *mockApp) ConfigRegister() interface{}                { return nil }
func (m *mockApp) CronRegistry() interface{}                  { return nil }
func (m *mockApp) GrpcRegistry() interface{}                  { return nil }
func (m *mockApp) GrpcConnRegistry() interface{}              { return nil }
func (m *mockApp) JwtRegister() interface{}                   { return nil }
func (m *mockApp) LanguageRegister() interface{}              { return nil }
func (m *mockApp) LockerRegistry() interface{}                { return nil }
func (m *mockApp) MemoryTaskRegister() interface{}            { return nil }
func (m *mockApp) MetricsRegister() interface{}               { return nil }
func (m *mockApp) MongoRegister() interface{}                 { return nil }
func (m *mockApp) QueueRegistry() interface{}                 { return nil }
func (m *mockApp) RabbitTaskRegister() interface{}            { return nil }
func (m *mockApp) RocketTaskRegister() interface{}            { return nil }
func (m *mockApp) RiakRegister() interface{}                  { return nil }
func (m *mockApp) ServerRegistry() interface{}                { return nil }
func (m *mockApp) SenderRegistry() interface{}                { return nil }
func (m *mockApp) StateRegistry() interface{}                 { return nil }
func (m *mockApp) TaskRegister() interface{}                  { return nil }
func (m *mockApp) TusUploaderRegister() interface{}           { return nil }
func (m *mockApp) WebSocketRegister() interface{}             { return nil }
func (m *mockApp) Lang(context.Context, string) string        { return "" }
func (m *mockApp) Config(context.Context, string) interface{} { return nil }
func (m *mockApp) GetQueueMessage(id, routingKey string, value map[string]interface{}) (interface{}, error) {
	return nil, nil
}
func (m *mockApp) VerifyRegistry() interface{} { return nil }

// mock boots and runners
func mockInitializeOk() boot.Initialize {
	return &testInit{"ok", false}
}
func mockInitializeFail() boot.Initialize {
	return &testInit{"fail", true}
}

type testInit struct {
	name string
	fail bool
}

func (t *testInit) String() string { return t.name }
func (t *testInit) Init(context.Context, app.IRuntime) error {
	if t.fail {
		return errors.New("init fail")
	}
	return nil
}

func TestBootstrap_Lifecycle(t *testing.T) {
	var beforeCalled, afterCalled int32
	var runnerRun int32
	b := NewBootstrap(context.Background(), &mockApp{})
	b.Before(func(context.Context, app.IRuntime) error {
		atomic.AddInt32(&beforeCalled, 1)
		return nil
	}).Process(mockInitializeOk()).Runner(func(ctx context.Context, app app.IRuntime) error {
		atomic.AddInt32(&runnerRun, 1)
		return nil
	}).After(func(context.Context, app.IRuntime) error {
		atomic.AddInt32(&afterCalled, 1)
		return nil
	})

	go func() {
		// 模拟发信号退出
		time.Sleep(200 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if beforeCalled != 1 || afterCalled != 1 || runnerRun != 1 {
		t.Fatalf("phase not executed correctly: before=%d, run=%d, after=%d", beforeCalled, runnerRun, afterCalled)
	}
}

func TestBootstrap_InitFail(t *testing.T) {
	b := NewBootstrap(context.Background(), &mockApp{})
	b.Process(mockInitializeFail())
	if err := b.Run(); err == nil {
		t.Fatal("expected error on init fail")
	}
}

func TestBootstrap_HealthCheck(t *testing.T) {
	b := NewBootstrap(context.Background(), &mockApp{})
	var hcall int32
	b.Process(mockInitializeOk()).WithHealthCheck(func(ctx context.Context) error {
		atomic.AddInt32(&hcall, 1)
		return nil
	})
	go func() {
		time.Sleep(200 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hcall != 1 {
		t.Fatal("health check not called")
	}
}

func TestBootstrap_RunnerAutoRestart(t *testing.T) {
	var runCount int32
	b := NewBootstrap(context.Background(), &mockApp{}).WithRunnerErrorPolicy("restart")
	b.Runner(func(ctx context.Context, a app.IRuntime) error {
		atomic.AddInt32(&runCount, 1)
		if runCount < 2 {
			return errors.New("fail to trigger restart")
		}
		// 模拟持久阻塞
		select {
		case <-ctx.Done():
			return nil
		}
	})
	go func() {
		time.Sleep(800 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if runCount < 2 {
		t.Fatal("restart policy not triggered")
	}
}

func TestBootstrap_HotReloadSignal(t *testing.T) {
	b := NewBootstrap(context.Background(), &mockApp{})
	b.Runner(func(ctx context.Context, a app.IRuntime) error {
		select {
		case <-ctx.Done():
			return nil
		}
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGHUP)
		// 再sleep等待reload后发SIGTERM退出
		time.Sleep(200 * time.Millisecond)
		p.Signal(syscall.SIGTERM)
	}()
	if err := b.Run(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}
