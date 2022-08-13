package ws

type Instance struct {
	ServerId    *uint64      `json:"serverId"`
	ConnManager *ConnManager `json:"connMgr"`
	Merge       *Merger      `json:"merger"`
	Cfg         *Config      `json:"config"`
	Stat        *Statistics  `json:"stats"`
}

func New() *Instance {
	return &Instance{}
}

func (i *Instance) GetInstance() *Instance {
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
