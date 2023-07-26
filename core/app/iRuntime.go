package app

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/plus-core/core/cache"
	"github.com/jxo-me/plus-core/core/locker"
	"github.com/jxo-me/plus-core/core/queue"
	reg "github.com/jxo-me/plus-core/core/registry"
)

type IRuntime interface {
	ServerRegistry() reg.IRegistry[*ghttp.Server]
	QueueRegistry() reg.IRegistry[queue.IQueue]
	CacheRegistry() reg.IRegistry[cache.ICache]
	LockerRegistry() reg.IRegistry[locker.ILocker]
}
