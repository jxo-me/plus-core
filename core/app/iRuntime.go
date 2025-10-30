package app

import (
	"context"

	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/jxo-me/plus-core/core/v2/cache"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	"github.com/jxo-me/plus-core/core/v2/queue"
	reg "github.com/jxo-me/plus-core/core/v2/registry"
	"github.com/jxo-me/plus-core/core/v2/task"
	"github.com/jxo-me/plus-core/pkg/v2/security/auth"
)

// ICacher    - 缓存能力运行时接口，聚合所有cache类型自动注册/热加载。
// IQueuer    - 消息队列能力接口，统一多种MQ注册与操作。
// ITasker    - 任务调度/处理能力抽象
// ISecurity  - 认证、校验相关接口。
// ILifeCycle - 生命周期批量管理。
// IPlatformRuntime - 能力总线聚合生命周期治理，可用于企业级平台统一控制。

// ICacher 包含所有缓存能力的动态注册和获取接口。
type ICacher interface {
	// CacheRegistry 返回缓存注册中心（支持多种类型cache能力的动态按需获取与注入）。
	CacheRegistry() reg.IRegistry[cache.ICache]
}

// IQueuer 消息队列注册统一接口。
type IQueuer interface {
	// QueueRegistry 提供所有队列能力的总注册接口。
	QueueRegistry() reg.IRegistry[queue.IQueue]
}

// ITasker 任务能力（含内存、各种调度服务）统一管理抽象接口。
type ITasker interface {
	// MemoryTaskRegister 返回内存任务服务注册。
	MemoryTaskRegister() reg.IRegistry[task.MemoryService]
	// TaskRegister 通用任务调度注册。
	TaskRegister() reg.IRegistry[task.TasksService]
}

// ISecurity 聚合认证、鉴权、校验注册能力。
type ISecurity interface {
	// JwtRegister JWT鉴权注册中心。
	JwtRegister() reg.IRegistry[*jwt.GfJWTMiddleware]
	// VerifyRegistry 签名/校验器注册。
	VerifyRegistry() reg.IRegistry[*auth.Verifier]
}

// ILifeCycle 支持批量能力的启动/停/健康批量管理。
type ILifeCycle interface {
	// StartAll 启动全部注册能力
	StartAll() error
	// StopAll 停止全部模块
	StopAll() error
	// CheckHealth 健康检查，包括自愈恢复探测。
	CheckHealth() error
}

// IPlatformRuntime 总线接口+生命周期能力。
type IPlatformRuntime interface {
	IRuntime
	ILifeCycle
}

// Reloadable 能力可选接口：用于在配置变更时执行热重载
// 实现方应保证幂等与线程安全

type Reloadable interface {
	Reload(ctx context.Context) error
}

// Configurable 能力可选接口：用于在配置变更时接收新配置
// 建议与 Reload 结合，先 ApplyConfig 再 Reload 生效

type Configurable interface {
	ApplyConfig(cfg interface{}) error
}

// 企业统一能力总线（分域能力接口聚合）
//
//go:generate go doc -all
type IRuntime interface {
	ICacher
	IQueuer
	ITasker
	ISecurity
	// ... other domain interfaces ...
	// 兼容原接口，便于平滑迁移
	Lang(ctx context.Context, langKey string) string
	Config(ctx context.Context, key string) *gvar.Var
	GetQueueMessage(id, routingKey string, value map[string]interface{}) (messageLib.IMessage, error)
}
