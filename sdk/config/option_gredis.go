package config

import (
	"context"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	insGRedis = cGRedis{}
)

type cGRedis struct {
	Client *gredis.Redis
}

type GRedisOptions struct {
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	Tls      *Tls   `yaml:"tls" json:"tls"`
}

// GRedis 获取redis客户端
func GRedis() *cGRedis {
	return &insGRedis
}

// GetClient 获取redis客户端
func (c *cGRedis) GetClient() *gredis.Redis {
	return c.Client
}

// SetClient 设置redis客户端
func (c *cGRedis) SetClient(ctx context.Context, r *gredis.Redis) *cGRedis {
	if c.Client != nil && c.Client != r {
		err := c.Client.Close(ctx)
		if err != nil {
			glog.Warning(ctx, "cGRedis close error:", err)
		}
	}
	c.Client = r
	return c
}

func (e GRedisOptions) GetClientOptions(ctx context.Context, s *Settings) (*gredis.Config, error) {
	address, err := s.Cfg().Get(ctx, "redis.default.address", "")
	if err != nil {
		return nil, err
	}
	db, err := s.Cfg().Get(ctx, "redis.default.db", "0")
	if err != nil {
		return nil, err
	}
	pass, err := s.Cfg().Get(ctx, "redis.default.pass", "")
	if err != nil {
		return nil, err
	}
	r := &gredis.Config{
		Address: address.String(),
		Pass:    pass.String(),
		Db:      db.Int(),
	}
	r.TLSConfig, err = getTLS(e.Tls)
	return r, err
}
