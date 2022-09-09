package ws

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"sync/atomic"
)

var insInstance = Instance{}

type Instance struct {
	ServerId    *uint64      `json:"serverId"`
	ConnManager *ConnManager `json:"connMgr"`
	Merge       *Merger      `json:"merger"`
	Cfg         *Config      `json:"config"`
	Stat        *Statistics  `json:"stats"`
}

func NewSocket(id *uint64, cfg *Config) *Instance {
	// InitStats
	insInstance.Stat = GetStats()
	// InitConnMgr
	insInstance.ConnManager = NewConnMgr(cfg)
	// InitMerger
	insInstance.Merge = NewMerger(cfg)
	// serverId
	insInstance.ServerId = id
	// config
	insInstance.Cfg = cfg

	return &insInstance
}

func (i *Instance) GetInstance() *Instance {
	return i
}

func (i *Instance) Connection(ctx context.Context, wsSocket *ghttp.WebSocket, routers *map[string]Service) *Instance {
	conn := NewConnection(ctx, atomic.AddUint64(i.ServerId, 1),
		wsSocket, i.Cfg.HeartbeatInterval,
		i.Cfg.InChannelSize,
		i.Cfg.OutChannelSize,
	)
	// 连接加入管理器, 可以推送端查找到
	i.ConnManager.AddConn(conn)
	conn.RouterHandle(ctx, i, routers)
	return i
}
