# sdk/config 指南

## 生命周期与启动器 Bootstrap
- Before/Process/Runner/After 四阶段
- 统一信号处理：SIGINT/SIGTERM 退出，SIGHUP 触发热重载
- 扩展：
  - Watchers(…): 注册外部配置变更监听
  - OnReload(…): 重载完成后的横切钩子
  - EnableAutoReloadCapabilities(): 自动对实现能力接口的对象执行ApplyConfig→Reload
  - EnableReloadDebounce(d): 重载去抖，合并高频事件
  - WithReloadBuffer(n): 调整重载事件缓冲

## 能力对象接口（可选实现）
- ApplyConfig(cfg interface{}) error：接收最新配置（由 Provider 提供）
- Reload(ctx context.Context) error：基于已应用的配置热重载
- 建议顺序：先 ApplyConfig → 后 Reload

## 能力注册中心（插件）
- 通过 `WithAbilityRegistry(view)` 注入能力注册视图（包含 Get/List）。
- 重载时会对插件执行 ApplyConfig → Reload。
- 热插拔：
  - Load(name) 优先调用插件 `Start(ctx)`（若实现），否则执行自定义 Loader。
  - Unload(name) 优先调用插件 `Stop(ctx)`（若实现），否则执行自定义 Unloader。

## 配置提供器与监听器
- WithConfigProvider(fn)：提供最新配置（任意结构，常为 map[string]any / struct）。
- Watchers(fn)：外部监听器检测到配置变更后调用 onChange()，统一通过 handleReload→worker 串行与去抖处理。

## Consul KV Bridge（示例，无外部依赖）
文件：`consul_kv_bridge.go`
- NewConsulKVBridge(addr, key, token, wait)
  - 基于 HTTP 阻塞查询监听 KV 变化
  - 支持 Value 为 base64 或明文 JSON，解析为 `map[string]any`
- Provider(ctx)：返回最新配置
- Watcher(ctx,onChange)：检测变更后调用 onChange

### Consul 接线示例
```go
ctx := context.Background()

b := config.NewBootstrap(ctx, myRuntime).
    EnableAutoReloadCapabilities().
    EnableReloadDebounce(200*time.Millisecond)

br := config.NewConsulKVBridge(
    "http://127.0.0.1:8500",
    "config/app.json",
    "",                 // token 可为空
    60*time.Second,      // 阻塞等待
)

b.WithConfigProvider(br.Provider).
  Watchers(br.Watcher)

b.OnReload(func(ctx context.Context, app app.IRuntime) error {
    // 横切逻辑：例如刷新路由/统计/缓存镜像等
    return nil
})

if err := b.Run(); err != nil { panic(err) }
```

### 能力对象实现示例
```go
type MyCap struct {
    cfg map[string]any
}

func (c *MyCap) ApplyConfig(cfg interface{}) error {
    if m, ok := cfg.(map[string]any); ok {
        c.cfg = m
    }
    return nil
}

func (c *MyCap) Reload(ctx context.Context) error {
    // 根据 c.cfg 重建连接池/更新路由/刷新策略等
    return nil
}
```

## 测试
- `bootstrap_reload_test.go`：watcher/SIGHUP重载
- `bootstrap_apply_reload_test.go`：ApplyConfig→Reload 链路
- `bootstrap_debounce_test.go`：去抖合并
- `bootstrap_ability_registry_test.go`：AbilityRegistry 插件 Start/Stop/ApplyConfig/Reload
- `consul_kv_bridge_test.go`：base64/明文 JSON 解析

## 建议
- Provider 返回带版本号/校验和的配置，processReload 做版本去重。
- 批量 Apply 成功后再统一 Reload，必要时设计回滚策略。
- 为重载流程与能力重载加埋点与日志，便于观测与告警。

## 最小可运行示例 main.go
```go
package main

import (
	"context"
	"time"

	cfg "github.com/jxo-me/plus-core/sdk/v2/config"
)

type myRuntime struct{}

func main() {
	ctx := context.Background()
	b := cfg.NewBootstrap(ctx, &myRuntime{}).
		EnableAutoReloadCapabilities().
		EnableReloadDebounce(200 * time.Millisecond)

	br := cfg.NewConsulKVBridge(
		"http://127.0.0.1:8500",
		"config/app.json",
		"",                // 可填token
		60*time.Second,
	)

	b.WithConfigProvider(br.Provider).
		Watchers(br.Watcher).
		OnReload(func(ctx context.Context, _ any) error {
			// 可增加横切逻辑，如刷新路由
			return nil
		})

	if err := b.Run(); err != nil { panic(err) }
}

## Consul / etcd 配置规范说明
- KV 建议：以 JSON 文本存储整体配置（或分模块JSON），示例：
  - key: `config/app.json`
  - value(JSON): `{"db":{"dsn":"..."},"cache":{"ttl":60}}`
- Consul：
  - Bridge 已支持 Value base64 和明文 JSON 两种形式。
  - 建议开启 ACL 使用 token；可配合 Namespace/Partition 按环境区隔。
- etcd：
  - 推荐将 etcd watcher 以与 Consul Bridge 类似方式实现：阻塞 watch + Provider 返回最新合并配置。
  - 将多 key 合并为一份配置对象（如以前缀聚合），在 Provider 中统一返回。
- 版本与幂等：
  - 建议在 JSON 中加入 `version` 或 `checksum` 字段，`WithConfigProvider` 返回后在 `processReload()` 做版本去重。
- 安全与审计：
  - 配置敏感字段建议加密存储（KMS/密钥服务），在 Provider 解密后再下发给能力层。
- 变更流程：
  - KV 更新 → Watcher onChange → Bootstrap 重载队列（去抖合并）→ ApplyConfig → Reload → OnReload 横切。
