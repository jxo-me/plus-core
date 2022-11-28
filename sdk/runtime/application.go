package runtime

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/gf-metrics"
	"github.com/jxo-me/gfbot"
	"github.com/jxo-me/plus-core/sdk/config"
	"github.com/jxo-me/plus-core/sdk/cron"
	"github.com/jxo-me/plus-core/sdk/pkg/tus"
	"github.com/jxo-me/plus-core/sdk/pkg/ws"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
	"github.com/jxo-me/plus-core/sdk/task"
	"sync"
)

type Application struct {
	server          *ghttp.Server
	casbins         map[string]*casbin.SyncedEnforcer
	mux             sync.RWMutex
	jwt             map[string]*jwt.GfJWTMiddleware
	lang            *gi18n.Manager
	config          *gcfg.Config
	settings        *config.Settings
	cache           storage.AdapterCache
	locker          storage.AdapterLocker
	websocket       *ws.Instance
	crontab         cron.Adapter
	taskService     task.TasksService
	rabbitmqService task.RabbitMqService
	rocketMqService task.RocketMqService
	memoryService   task.MemoryService
	queue           map[string]storage.AdapterQueue
	tus             *tus.Uploader
	monitor         *metrics.Monitor
	bot             *telebot.Bot
	botHook         *telebot.Hook
}

// NewConfig 默认值
func NewConfig() *Application {
	return &Application{
		casbins: make(map[string]*casbin.SyncedEnforcer),
		jwt:     make(map[string]*jwt.GfJWTMiddleware),
		queue:   make(map[string]storage.AdapterQueue),
	}
}

func (e *Application) SetCron(srv cron.Adapter) {
	e.crontab = srv
}

func (e *Application) Cron() cron.Adapter {
	return e.crontab
}

func (e *Application) SetTask(srv task.TasksService) {
	e.taskService = srv
}

func (e *Application) Task() task.TasksService {
	return e.taskService
}

func (e *Application) SetTus(t *tus.Uploader) {
	e.tus = t
}

func (e *Application) Tus() *tus.Uploader {
	return e.tus
}

func (e *Application) SetRabbitTask(srv task.RabbitMqService) {
	e.rabbitmqService = srv
}

func (e *Application) RabbitTask() task.RabbitMqService {
	return e.rabbitmqService
}

func (e *Application) SetRocketMqTask(srv task.RocketMqService) {
	e.rocketMqService = srv
}

func (e *Application) RocketMqTask() task.RocketMqService {
	return e.rocketMqService
}

func (e *Application) SetMemoryTask(srv task.MemoryService) {
	e.memoryService = srv
}

func (e *Application) MemoryTask() task.MemoryService {
	return e.memoryService
}

func (e *Application) SetServer(srv *ghttp.Server) {
	e.server = srv
}

func (e *Application) Server() *ghttp.Server {
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
		glog.Error(ctx, "Runtime Get Global Config error:", err)
	}
	return c
}

func (e *Application) SetSettings(c *config.Settings) {
	e.settings = c
}

func (e *Application) Settings() *config.Settings {
	return e.settings
}

// SetCache 设置缓存
func (e *Application) SetCache(c storage.AdapterCache) {
	e.cache = c
}

func (e *Application) Cache() storage.AdapterCache {
	return NewCache("", e.cache)
}

// GetCacheAdapter 获取缓存
func (e *Application) GetCacheAdapter() storage.AdapterCache {
	return NewCache("", e.cache)
}

// GetCachePrefix 获取带租户标记的cache
func (e *Application) GetCachePrefix(key string) storage.AdapterCache {
	return NewCache(key, e.cache)
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

func (e *Application) GetMemoryQueue(prefix string) storage.AdapterQueue {
	e.mux.Lock()
	defer e.mux.Unlock()
	if q, ok := e.queue[fmt.Sprintf("%s_%s", prefix, config.MemoryQueueName)]; ok {
		return q
	}
	return nil
}

func (e *Application) GetRabbitQueue(prefix string) storage.AdapterQueue {
	e.mux.Lock()
	defer e.mux.Unlock()
	if q, ok := e.queue[fmt.Sprintf("%s_%s", prefix, config.RabbitmqQueueName)]; ok {
		return q
	}
	return nil
}

func (e *Application) GetRocketQueue(prefix string) storage.AdapterQueue {
	e.mux.Lock()
	defer e.mux.Unlock()
	if q, ok := e.queue[fmt.Sprintf("%s_%s", prefix, config.RocketQueueName)]; ok {
		return q
	}
	return nil
}

// SetQueueAdapter 设置队列适配器
func (e *Application) SetQueueAdapter(key string, c storage.AdapterQueue) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.queue[fmt.Sprintf("%s_%s", key, c.String())] = NewQueue(key, c)
}

// GetQueueAdapter 获取队列适配器
func (e *Application) GetQueueAdapter(key string) storage.AdapterQueue {
	e.mux.Lock()
	defer e.mux.Unlock()
	// 优先返回全局
	if j, ok := e.queue[fmt.Sprintf("*_%s", key)]; ok {
		return j
	}
	return e.queue[key]
}

// GetQueueMessage 获取队列需要用的message
func (e *Application) GetQueueMessage(id, routingKey string, value map[string]interface{}) (storage.Messager, error) {
	message := &queue.Message{}
	message.SetId(id)
	message.SetRoutingKey(routingKey)
	message.SetValues(value)
	return message, nil
}

// SetLockerAdapter 设置分布式锁
func (e *Application) SetLockerAdapter(c storage.AdapterLocker) {
	e.locker = c
}

// GetLockerAdapter 获取分布式锁
func (e *Application) GetLockerAdapter() storage.AdapterLocker {
	return NewLocker("", e.locker)
}

func (e *Application) GetLockerPrefix(key string) storage.AdapterLocker {
	return NewLocker(key, e.locker)
}

func (e *Application) SetMetrics(m *metrics.Monitor) {
	e.monitor = m
}

func (e *Application) Monitor() *metrics.Monitor {
	return e.monitor
}

func (e *Application) SetBot(b *telebot.Bot) {
	e.bot = b
}

func (e *Application) Bot() *telebot.Bot {
	return e.bot
}

func (e *Application) SetBotHook(b *telebot.Hook) {
	e.botHook = b
}

func (e *Application) BotHook() *telebot.Hook {
	return e.botHook
}
