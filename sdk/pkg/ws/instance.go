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
	insInstance.Stat = Stats()
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
	conn := NewConnection(ctx, atomic.AddUint64(i.Sid(), 1),
		wsSocket, i.Config().HeartbeatInterval,
		i.Config().InChannelSize,
		i.Config().OutChannelSize,
	)
	// 连接加入管理器, 可以推送端查找到
	i.ConnMgr().AddConn(conn)
	conn.RouterHandle(ctx, i, routers)
	return i
}

func (i *Instance) Sid() *uint64 {
	return i.ServerId
}

func (i *Instance) SetServerId(id uint64) *Instance {
	i.ServerId = &id
	return i
}

func (i *Instance) GetServerId() *uint64 {
	return i.ServerId
}

func (i *Instance) ConnMgr() *ConnManager {
	return i.ConnManager
}

func (i *Instance) SetConnMgr(c *ConnManager) *Instance {
	i.ConnManager = c
	return i
}

func (i *Instance) GetConnMgr() *ConnManager {
	return i.ConnManager
}

func (i *Instance) Config() *Config {
	return i.Cfg
}

func (i *Instance) SetConfig(c *Config) *Instance {
	i.Cfg = c
	return i
}

func (i *Instance) GetConfig() *Config {
	return i.Cfg
}

func (i *Instance) Merger() *Merger {
	return i.Merge
}

func (i *Instance) SetMerger(m *Merger) *Instance {
	i.Merge = m
	return i
}

func (i *Instance) GetMerger() *Merger {
	return i.Merge
}

func (i *Instance) Stats() *Statistics {
	return i.Stat
}

func (i *Instance) SetStats(s *Statistics) *Instance {
	i.Stat = s
	return i
}

func (i *Instance) GetStats() *Statistics {
	return i.Stat
}
