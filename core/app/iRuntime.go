package app

import (
	"context"
	"github.com/casbin/casbin/v2"
	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcfg"
	metrics "github.com/jxo-me/gf-metrics"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/plus-core/core/cache"
	"github.com/jxo-me/plus-core/core/cron"
	"github.com/jxo-me/plus-core/core/locker"
	"github.com/jxo-me/plus-core/core/message"
	"github.com/jxo-me/plus-core/core/queue"
	reg "github.com/jxo-me/plus-core/core/registry"
	"github.com/jxo-me/plus-core/core/server"
	"github.com/jxo-me/plus-core/core/task"
	"github.com/jxo-me/plus-core/pkg/tus"
	"github.com/jxo-me/plus-core/pkg/ws"
)

type IRuntime interface {
	ServerRegistry() reg.IRegistry[server.IServer]
	QueueRegistry() reg.IRegistry[queue.IQueue]
	GetQueueMessage(id, routingKey string, value map[string]interface{}) (message.IMessage, error)
	CacheRegistry() reg.IRegistry[cache.ICache]
	LockerRegistry() reg.IRegistry[locker.ILocker]
	CronRegistry() reg.IRegistry[cron.ICron]
	BotRegistry() reg.IRegistry[*telebot.Bot]
	CasBinRegistry() reg.IRegistry[*casbin.SyncedEnforcer]
	JwtRegister() reg.IRegistry[*jwt.GfJWTMiddleware]
	LangRegister() reg.IRegistry[*gi18n.Manager]
	Lang(ctx context.Context, langKey string) string
	MetricsRegister() reg.IRegistry[*metrics.Monitor]
	ConfigRegister() reg.IRegistry[*gcfg.Config]
	Config(ctx context.Context, pattern string) *gvar.Var
	TaskRegister() reg.IRegistry[task.TasksService]
	RabbitTaskRegister() reg.IRegistry[task.RabbitMqService]
	RocketTaskRegister() reg.IRegistry[task.RocketMqService]
	MemoryTaskRegister() reg.IRegistry[task.MemoryService]
	TusUploaderRegister() reg.IRegistry[*tus.Uploader]
	WebSocketRegister() reg.IRegistry[*ws.Instance]
}
