package runtime

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/jxo-me/plus-core/sdk/config"
	"github.com/jxo-me/plus-core/sdk/cron"
	"github.com/jxo-me/plus-core/sdk/pkg/ws"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/task"
)

type Runtime interface {
	// SetServer Http Server
	SetServer(srv *ghttp.Server)
	Server() *ghttp.Server

	// SetCasbin casbin module
	SetCasbin(key string, enforcer *casbin.SyncedEnforcer)
	GetCasbin() map[string]*casbin.SyncedEnforcer
	GetCasbinKey(key string) *casbin.SyncedEnforcer
	// SetJwt jwt module
	SetJwt(key string, jwtIns *jwt.GfJWTMiddleware)
	GetJwt() map[string]*jwt.GfJWTMiddleware
	GetJwtKey(moduleKey string) *jwt.GfJWTMiddleware
	// SetLang gi18n
	SetLang(lang *gi18n.Manager)
	GetLang() *gi18n.Manager
	Lang(ctx context.Context, langKey string) string
	// SetConfig config
	SetConfig(c *gcfg.Config)
	GetConfig() *gcfg.Config
	Config(ctx context.Context, pattern string) *gvar.Var
	SetSettings(c *config.Settings)
	Settings() *config.Settings
	// SetCacheAdapter cache
	SetCacheAdapter(storage.AdapterCache)
	Cache() storage.AdapterCache
	GetCacheAdapter() storage.AdapterCache
	GetCachePrefix(string) storage.AdapterCache
	// SetWebSocket websocket
	SetWebSocket(s *ws.Instance)
	WebSocket() *ws.Instance
	GetWebSocket() *ws.Instance

	GetMemoryQueue(prefix string) storage.AdapterQueue
	GetRabbitQueue(prefix string) storage.AdapterQueue
	GetRocketQueue(prefix string) storage.AdapterQueue

	SetQueueAdapter(string, storage.AdapterQueue)
	GetQueueAdapter(key string) storage.AdapterQueue
	GetQueuePrefix(key string) storage.AdapterQueue
	GetQueueMessage(id, routingKey string, value map[string]interface{}) (storage.Messager, error)

	SetLockerAdapter(storage.AdapterLocker)
	GetLockerAdapter() storage.AdapterLocker
	GetLockerPrefix(string) storage.AdapterLocker

	SetCron(srv cron.Adapter)
	Cron() cron.Adapter

	SetTask(srv task.TasksService)
	Task() task.TasksService
	SetRabbitTask(srv task.RabbitMqService)
	RabbitTask() task.RabbitMqService
	SetRocketMqTask(srv task.RocketMqService)
	RocketMqTask() task.RocketMqService
	SetMemoryTask(srv task.MemoryService)
	MemoryTask() task.MemoryService
}
