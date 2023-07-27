package app

import (
	"github.com/casbin/casbin/v2"
	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	metrics "github.com/jxo-me/gf-metrics"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/plus-core/core/cache"
	"github.com/jxo-me/plus-core/core/cron"
	"github.com/jxo-me/plus-core/core/locker"
	"github.com/jxo-me/plus-core/core/queue"
	reg "github.com/jxo-me/plus-core/core/registry"
	"github.com/jxo-me/plus-core/core/task"
	"github.com/jxo-me/plus-core/pkg/tus"
	"github.com/jxo-me/plus-core/pkg/ws"
)

type IRuntime interface {
	BotRegistry() reg.IRegistry[*telebot.Bot]
	CacheRegistry() reg.IRegistry[cache.ICache]
	CasBinRegistry() reg.IRegistry[*casbin.SyncedEnforcer]
	ConfigRegister() reg.IRegistry[*gcfg.Config]
	CronRegistry() reg.IRegistry[cron.ICron]
	JwtRegister() reg.IRegistry[*jwt.GfJWTMiddleware]
	LanguageRegister() reg.IRegistry[*gi18n.Manager]
	LockerRegistry() reg.IRegistry[locker.ILocker]
	MemoryTaskRegister() reg.IRegistry[task.MemoryService]
	MetricsRegister() reg.IRegistry[*metrics.Monitor]
	QueueRegistry() reg.IRegistry[queue.IQueue]
	RabbitTaskRegister() reg.IRegistry[task.RabbitMqService]
	RocketTaskRegister() reg.IRegistry[task.RocketMqService]
	ServerRegistry() reg.IRegistry[*ghttp.Server]
	TaskRegister() reg.IRegistry[task.TasksService]
	TusUploaderRegister() reg.IRegistry[*tus.Uploader]
	WebSocketRegister() reg.IRegistry[*ws.Instance]
}
