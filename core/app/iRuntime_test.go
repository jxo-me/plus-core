package app

import (
	"context"
	"testing"

	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/jxo-me/plus-core/core/v2/cache"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	"github.com/jxo-me/plus-core/core/v2/queue"
	registry "github.com/jxo-me/plus-core/core/v2/registry"
	"github.com/jxo-me/plus-core/core/v2/task"
	"github.com/jxo-me/plus-core/pkg/v2/security/auth"
)

type mockCacheReg struct {
	data map[string]cache.ICache
}

func (m *mockCacheReg) Register(name string, v cache.ICache) error {
	if m.data == nil {
		m.data = make(map[string]cache.ICache)
	}
	m.data[name] = v
	return nil
}
func (m *mockCacheReg) Get(name string) cache.ICache    { return m.data[name] }
func (m *mockCacheReg) GetAll() map[string]cache.ICache { return m.data }
func (m *mockCacheReg) Unregister(name string)          { delete(m.data, name) }
func (m *mockCacheReg) IsRegistered(name string) bool   { _, ok := m.data[name]; return ok }

type mockQueueReg struct {
	data map[string]queue.IQueue
}

func (m *mockQueueReg) Register(name string, v queue.IQueue) error {
	if m.data == nil {
		m.data = make(map[string]queue.IQueue)
	}
	m.data[name] = v
	return nil
}
func (m *mockQueueReg) Get(name string) queue.IQueue    { return m.data[name] }
func (m *mockQueueReg) GetAll() map[string]queue.IQueue { return m.data }
func (m *mockQueueReg) Unregister(name string)          { delete(m.data, name) }
func (m *mockQueueReg) IsRegistered(name string) bool   { _, ok := m.data[name]; return ok }

type mockTaskReg struct {
	data map[string]task.TasksService
}

func (m *mockTaskReg) Register(name string, v task.TasksService) error {
	if m.data == nil {
		m.data = make(map[string]task.TasksService)
	}
	m.data[name] = v
	return nil
}
func (m *mockTaskReg) Get(name string) task.TasksService    { return m.data[name] }
func (m *mockTaskReg) GetAll() map[string]task.TasksService { return m.data }
func (m *mockTaskReg) Unregister(name string)               { delete(m.data, name) }
func (m *mockTaskReg) IsRegistered(name string) bool        { _, ok := m.data[name]; return ok }

type mockMemoryTaskReg struct {
	data map[string]task.MemoryService
}

func (m *mockMemoryTaskReg) Register(name string, v task.MemoryService) error {
	if m.data == nil {
		m.data = make(map[string]task.MemoryService)
	}
	m.data[name] = v
	return nil
}
func (m *mockMemoryTaskReg) Get(name string) task.MemoryService    { return m.data[name] }
func (m *mockMemoryTaskReg) GetAll() map[string]task.MemoryService { return m.data }
func (m *mockMemoryTaskReg) Unregister(name string)                { delete(m.data, name) }
func (m *mockMemoryTaskReg) IsRegistered(name string) bool         { _, ok := m.data[name]; return ok }

type mockJwtReg struct {
	data map[string]*jwt.GfJWTMiddleware
}

func (m *mockJwtReg) Register(name string, v *jwt.GfJWTMiddleware) error {
	if m.data == nil {
		m.data = make(map[string]*jwt.GfJWTMiddleware)
	}
	m.data[name] = v
	return nil
}
func (m *mockJwtReg) Get(name string) *jwt.GfJWTMiddleware    { return m.data[name] }
func (m *mockJwtReg) GetAll() map[string]*jwt.GfJWTMiddleware { return m.data }
func (m *mockJwtReg) Unregister(name string)                  { delete(m.data, name) }
func (m *mockJwtReg) IsRegistered(name string) bool           { _, ok := m.data[name]; return ok }

type mockVerifierReg struct {
	data map[string]*auth.Verifier
}

func (m *mockVerifierReg) Register(name string, v *auth.Verifier) error {
	if m.data == nil {
		m.data = make(map[string]*auth.Verifier)
	}
	m.data[name] = v
	return nil
}
func (m *mockVerifierReg) Get(name string) *auth.Verifier    { return m.data[name] }
func (m *mockVerifierReg) GetAll() map[string]*auth.Verifier { return m.data }
func (m *mockVerifierReg) Unregister(name string)            { delete(m.data, name) }
func (m *mockVerifierReg) IsRegistered(name string) bool     { _, ok := m.data[name]; return ok }

// mock实现
type demoRuntime struct{}

func (d *demoRuntime) CacheRegistry() registry.IRegistry[cache.ICache] {
	return &mockCacheReg{make(map[string]cache.ICache)}
}
func (d *demoRuntime) QueueRegistry() registry.IRegistry[queue.IQueue] {
	return &mockQueueReg{make(map[string]queue.IQueue)}
}
func (d *demoRuntime) MemoryTaskRegister() registry.IRegistry[task.MemoryService] {
	return &mockMemoryTaskReg{make(map[string]task.MemoryService)}
}
func (d *demoRuntime) TaskRegister() registry.IRegistry[task.TasksService] {
	return &mockTaskReg{make(map[string]task.TasksService)}
}
func (d *demoRuntime) JwtRegister() registry.IRegistry[*jwt.GfJWTMiddleware] {
	return &mockJwtReg{make(map[string]*jwt.GfJWTMiddleware)}
}
func (d *demoRuntime) VerifyRegistry() registry.IRegistry[*auth.Verifier] {
	return &mockVerifierReg{make(map[string]*auth.Verifier)}
}
func (d *demoRuntime) Lang(ctx context.Context, lang string) string     { return "zh" }
func (d *demoRuntime) Config(ctx context.Context, key string) *gvar.Var { return nil }
func (d *demoRuntime) GetQueueMessage(id, routingKey string, v map[string]interface{}) (messageLib.IMessage, error) {
	return nil, nil
}

// 生命周期能力mock
func (d *demoRuntime) StartAll() error    { return nil }
func (d *demoRuntime) StopAll() error     { return nil }
func (d *demoRuntime) CheckHealth() error { return nil }

func TestRuntime_CompositeContract(t *testing.T) {
	var _ ICacher = &demoRuntime{}
	var _ IQueuer = &demoRuntime{}
	var _ ITasker = &demoRuntime{}
	var _ ISecurity = &demoRuntime{}
	var _ ILifeCycle = &demoRuntime{}
	var _ IRuntime = &demoRuntime{}
	var _ IPlatformRuntime = &demoRuntime{}
}

func TestPlatform_LifeCycleCall(t *testing.T) {
	d := &demoRuntime{}
	if err := d.StartAll(); err != nil {
		t.Fatal("start fail")
	}
	if err := d.CheckHealth(); err != nil {
		t.Fatal("health fail")
	}
	if err := d.StopAll(); err != nil {
		t.Fatal("stop fail")
	}
}
