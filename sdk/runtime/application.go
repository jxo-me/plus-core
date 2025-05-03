package runtime

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
	metrics "github.com/jxo-me/gf-metrics"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/plus-core/core/v2/bucket"
	cacheLib "github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/jxo-me/plus-core/core/v2/cron"
	lockerLib "github.com/jxo-me/plus-core/core/v2/locker"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	reg "github.com/jxo-me/plus-core/core/v2/registry"
	"github.com/jxo-me/plus-core/core/v2/send"
	"github.com/jxo-me/plus-core/core/v2/task"
	"github.com/jxo-me/plus-core/pkg/v2/security/auth"
	"github.com/jxo-me/plus-core/pkg/v2/tus"
	"github.com/jxo-me/plus-core/pkg/v2/ws"
	"github.com/jxo-me/plus-core/sdk/v2/message"
	"github.com/jxo-me/plus-core/sdk/v2/registry"
	"github.com/lesismal/arpc"
	"github.com/zegl/goriak/v3"
	"google.golang.org/grpc"
)

type Application struct {
	arpcServerReg      reg.IRegistry[*arpc.Server]
	arpcClientReg      reg.IRegistry[*arpc.Client]
	botReg             reg.IRegistry[*telebot.Bot]
	cacheReg           reg.IRegistry[cacheLib.ICache]
	casBinReg          reg.IRegistry[*casbin.SyncedEnforcer]
	configReg          reg.IRegistry[*gcfg.Config]
	crontabReg         reg.IRegistry[cron.ICron]
	grpcReg            reg.IRegistry[*grpcx.GrpcServer]
	grpcConnReg        reg.IRegistry[*grpc.ClientConn]
	jwtReg             reg.IRegistry[*jwt.GfJWTMiddleware]
	languageReg        reg.IRegistry[*gi18n.Manager]
	lockerReg          reg.IRegistry[lockerLib.ILocker]
	memoryServiceReg   reg.IRegistry[task.MemoryService]
	metricsReg         reg.IRegistry[*metrics.Monitor]
	queueReg           reg.IRegistry[queueLib.IQueue]
	rabbitMqServiceReg reg.IRegistry[task.RabbitMqService]
	rocketMqServiceReg reg.IRegistry[task.RocketMqService]
	serverReg          reg.IRegistry[*ghttp.Server]
	senderReg          reg.IRegistry[send.ISender[send.ISendMsg]]
	stateReg           reg.IRegistry[bucket.IState]
	taskServiceReg     reg.IRegistry[task.TasksService]
	tusReg             reg.IRegistry[*tus.Uploader]
	websocketReg       reg.IRegistry[*ws.Instance]
	riakReg            reg.IRegistry[*goriak.Session]
	verifyReg          reg.IRegistry[*auth.Verifier]
}

// NewConfig 默认值
func NewConfig() *Application {
	return &Application{
		arpcServerReg:      new(registry.ArpcServerRegistry),
		arpcClientReg:      new(registry.ArpcClientRegistry),
		botReg:             new(registry.BotRegistry),
		cacheReg:           new(registry.CacheRegistry),
		casBinReg:          new(registry.CasBinRegistry),
		configReg:          new(registry.ConfigRegistry),
		crontabReg:         new(registry.CrontabRegistry),
		grpcReg:            new(registry.GrpcRegistry),
		grpcConnReg:        new(registry.GrpcConnRegistry),
		jwtReg:             new(registry.JwtRegistry),
		languageReg:        new(registry.LanguageRegistry),
		lockerReg:          new(registry.LockerRegistry),
		memoryServiceReg:   new(registry.MemoryServiceRegistry),
		metricsReg:         new(registry.MetricsRegistry),
		queueReg:           new(registry.QueueRegistry),
		rabbitMqServiceReg: new(registry.RabbitMqServiceRegistry),
		rocketMqServiceReg: new(registry.RocketMqServiceRegistry),
		serverReg:          new(registry.ServerRegistry),
		senderReg:          new(registry.SenderRegistry),
		stateReg:           new(registry.StateRegistry),
		taskServiceReg:     new(registry.TaskServiceRegistry),
		tusReg:             new(registry.TusRegistry),
		websocketReg:       new(registry.WebSocketRegistry),
		riakReg:            new(registry.RiakRegistry),
		verifyReg:          new(registry.VerifyRegistry),
	}
}

func (a *Application) ArpcServerRegistry() reg.IRegistry[*arpc.Server] {
	return a.arpcServerReg
}

func (a *Application) ArpcClientRegistry() reg.IRegistry[*arpc.Client] {
	return a.arpcClientReg
}

func (a *Application) BotRegistry() reg.IRegistry[*telebot.Bot] {
	return a.botReg
}

func (a *Application) CacheRegistry() reg.IRegistry[cacheLib.ICache] {
	return a.cacheReg
}

func (a *Application) CasBinRegistry() reg.IRegistry[*casbin.SyncedEnforcer] {
	return a.casBinReg
}

func (a *Application) ConfigRegister() reg.IRegistry[*gcfg.Config] {
	return a.configReg
}

func (a *Application) Config(ctx context.Context, key string) *gvar.Var {
	r, err := a.configReg.Get("").Get(ctx, key)
	if err != nil {
		glog.Error(ctx, "get global config error:", err.Error())
		return gvar.New(nil)
	}
	return r
}

func (a *Application) CronRegistry() reg.IRegistry[cron.ICron] {
	return a.crontabReg
}

func (a *Application) GrpcRegistry() reg.IRegistry[*grpcx.GrpcServer] {
	return a.grpcReg
}

func (a *Application) GrpcConnRegistry() reg.IRegistry[*grpc.ClientConn] {
	return a.grpcConnReg
}

func (a *Application) JwtRegister() reg.IRegistry[*jwt.GfJWTMiddleware] {
	return a.jwtReg
}

func (a *Application) LanguageRegister() reg.IRegistry[*gi18n.Manager] {
	return a.languageReg
}

func (a *Application) Lang(ctx context.Context, langKey string) string {
	return a.languageReg.Get("").Translate(ctx, fmt.Sprintf(`{#%s}`, langKey))
}

func (a *Application) LockerRegistry() reg.IRegistry[lockerLib.ILocker] {
	return a.lockerReg
}

func (a *Application) MemoryTaskRegister() reg.IRegistry[task.MemoryService] {
	return a.memoryServiceReg
}

func (a *Application) MetricsRegister() reg.IRegistry[*metrics.Monitor] {
	return a.metricsReg
}

func (a *Application) QueueRegistry() reg.IRegistry[queueLib.IQueue] {
	return a.queueReg
}

func (a *Application) RabbitTaskRegister() reg.IRegistry[task.RabbitMqService] {
	return a.rabbitMqServiceReg
}

func (a *Application) RocketTaskRegister() reg.IRegistry[task.RocketMqService] {
	return a.rocketMqServiceReg
}

func (a *Application) RiakRegister() reg.IRegistry[*goriak.Session] {
	return a.riakReg
}

func (a *Application) ServerRegistry() reg.IRegistry[*ghttp.Server] {
	return a.serverReg
}

func (a *Application) SenderRegistry() reg.IRegistry[send.ISender[send.ISendMsg]] {
	return a.senderReg
}

func (a *Application) StateRegistry() reg.IRegistry[bucket.IState] {
	return a.stateReg
}

func (a *Application) TaskRegister() reg.IRegistry[task.TasksService] {
	return a.taskServiceReg
}

func (a *Application) TusUploaderRegister() reg.IRegistry[*tus.Uploader] {
	return a.tusReg
}

func (a *Application) WebSocketRegister() reg.IRegistry[*ws.Instance] {
	return a.websocketReg
}

// GetQueueMessage 获取队列需要用的message
func (a *Application) GetQueueMessage(id, routingKey string, value map[string]interface{}) (messageLib.IMessage, error) {
	m := &message.Message{}
	m.SetId(id)
	m.SetRoutingKey(routingKey)
	m.SetValue(value)
	return m, nil
}

func (a *Application) VerifyRegistry() reg.IRegistry[*auth.Verifier] {
	return a.verifyReg
}
