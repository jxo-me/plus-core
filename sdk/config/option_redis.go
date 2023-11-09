package config

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/redis/go-redis/v9"
)

var insRedis = cRedis{}

type cRedis struct {
	Client *redis.Client
}

// Redis 获取redis客户端
func Redis() *cRedis {
	return &insRedis
}

// GetClient 获取redis客户端
func (c *cRedis) GetClient() *redis.Client {
	return c.Client
}

// SetClient 设置redis客户端
func (c *cRedis) SetClient(ctx context.Context, r *redis.Client) *cRedis {
	if c.Client != nil && c.Client != r {
		err := c.Client.Shutdown(ctx)
		if err != nil {
			glog.Warning(ctx, "cRedis Shutdown error:", err)
		}
	}
	c.Client = r
	return c
}

type RedisConnectOptions struct {
	Network    string `yaml:"network" json:"network"`
	Addr       string `yaml:"addr" json:"addr"`
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
	DB         int    `yaml:"db" json:"db"`
	PoolSize   int    `yaml:"pool_size" json:"pool_size"`
	Tls        *Tls   `yaml:"tls" json:"tls"`
	MaxRetries int    `yaml:"max_retries" json:"max_retries"`
}

func (e *RedisConnectOptions) GetRedisOptions(ctx context.Context, s *Settings) (*redis.Options, error) {
	address, err := s.Cfg().Get(ctx, "redis.default.address", "")
	if err != nil {
		return nil, err
	}
	db, err := s.Cfg().Get(ctx, "redis.default.db", "0")
	if err != nil {
		return nil, err
	}
	network, err := s.Cfg().Get(ctx, "redis.default.network", "tcp")
	if err != nil {
		return nil, err
	}
	pass, err := s.Cfg().Get(ctx, "redis.default.pass", "")
	if err != nil {
		return nil, err
	}
	username, err := s.Cfg().Get(ctx, "redis.default.username", "")
	if err != nil {
		return nil, err
	}
	poolSize, err := s.Cfg().Get(ctx, "redis.default.pool_size", "10")
	if err != nil {
		return nil, err
	}
	maxRetries, err := s.Cfg().Get(ctx, "redis.default.max_retries", "0")
	if err != nil {
		return nil, err
	}
	// cert key ca
	cert, err := s.Cfg().Get(ctx, "redis.default.cert", "")
	if err != nil {
		return nil, err
	}
	key, err := s.Cfg().Get(ctx, "redis.default.key", "")
	if err != nil {
		return nil, err
	}
	ca, err := s.Cfg().Get(ctx, "redis.default.ca", "")
	if err != nil {
		return nil, err
	}
	tls := &Tls{
		Cert: cert.String(),
		Ca:   ca.String(),
		Key:  key.String(),
	}
	r := &redis.Options{
		Network:    network.String(),
		Addr:       address.String(),
		Username:   username.String(),
		Password:   pass.String(),
		DB:         db.Int(),
		MaxRetries: maxRetries.Int(),
		PoolSize:   poolSize.Int(),
	}
	r.TLSConfig, err = getTLS(tls)
	return r, err
}
