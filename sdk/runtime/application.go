package runtime

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/pkg/amqp/pool"
	"github.com/jxo-me/plus-core/sdk/pkg/ws"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
	"sync"
)

type Application struct {
	server      *ghttp.Server
	casbins     map[string]*casbin.SyncedEnforcer
	mux         sync.RWMutex
	jwt         map[string]*jwt.GfJWTMiddleware
	lang        *gi18n.Manager
	config      *gcfg.Config
	cache       *gcache.Cache
	websocket   *ws.Instance
	memoryQueue storage.AdapterQueue
	amqp        map[string]*pool.ConnPool
	queue       storage.AdapterQueue
}

// NewConfig 默认值
func NewConfig() *Application {
	return &Application{
		casbins:     make(map[string]*casbin.SyncedEnforcer),
		jwt:         make(map[string]*jwt.GfJWTMiddleware),
		amqp:        make(map[string]*pool.ConnPool),
		memoryQueue: queue.NewMemory(10000),
	}
}

func (e *Application) SetServer(srv *ghttp.Server) {
	e.server = srv
}

func (e *Application) GetServer() *ghttp.Server {
	return e.server
}

func (e *Application) SetCasbin(key string, enforcer *casbin.SyncedEnforcer) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.casbins[key] = enforcer
}

func (e *Application) GetCasbin() map[string]*casbin.SyncedEnforcer {
	return e.casbins
}

// GetCasbinKey 根据key获取casBin
func (e *Application) GetCasbinKey(key string) *casbin.SyncedEnforcer {
	e.mux.Lock()
	defer e.mux.Unlock()
	if c, ok := e.casbins["*"]; ok {
		return c
	}
	return e.casbins[key]
}

func (e *Application) SetJwt(key string, jwtIns *jwt.GfJWTMiddleware) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.jwt[key] = jwtIns
}

func (e *Application) GetJwt() map[string]*jwt.GfJWTMiddleware {
	return e.jwt
}

// GetJwtKey 根据key获取Jwt
func (e *Application) GetJwtKey(moduleKey string) *jwt.GfJWTMiddleware {
	e.mux.Lock()
	defer e.mux.Unlock()
	if j, ok := e.jwt["*"]; ok {
		return j
	}
	return e.jwt[moduleKey]
}

// Lang 多语言翻译
func (e *Application) Lang(ctx context.Context, langKey string) string {
	return e.GetLang().Translate(ctx, fmt.Sprintf(`{#%s}`, langKey))
}

func (e *Application) SetLang(lang *gi18n.Manager) {
	e.lang = lang
}

func (e *Application) GetLang() *gi18n.Manager {
	return e.lang
}

func (e *Application) SetConfig(c *gcfg.Config) {
	e.config = c
}

func (e *Application) GetConfig() *gcfg.Config {
	return e.config
}

func (e *Application) Config(ctx context.Context, pattern string) *gvar.Var {
	c, err := e.GetConfig().Get(ctx, pattern)
	if err != nil {
		glog.Errorf(ctx, "Runtime Get Global Config error:", err)
	}
	return c
}

// Cache 获取缓存实例
func (e *Application) Cache() *gcache.Cache {
	return e.cache
}

// SetCache 设置缓存
func (e *Application) SetCache(c *gcache.Cache) {
	e.cache = c
}

// GetCache 获取缓存
func (e *Application) GetCache() *gcache.Cache {
	return e.cache
}

func (e *Application) SetWebSocket(s *ws.Instance) {
	e.websocket = s
}

func (e *Application) WebSocket() *ws.Instance {
	return e.websocket
}

func (e *Application) GetWebSocket() *ws.Instance {
	return e.websocket
}

func (e *Application) SetAmqp(key string, amqp *pool.ConnPool) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.amqp[key] = amqp
}

func (e *Application) GetAmqp() map[string]*pool.ConnPool {
	return e.amqp
}

// GetAmqpKey 根据key获取amqp
func (e *Application) GetAmqpKey(key string) *pool.ConnPool {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e, ok := e.amqp["*"]; ok {
		return e
	}
	return e.amqp[key]
}
