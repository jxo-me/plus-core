package config

import (
	"github.com/go-redis/redis/v7"
)

var (
	_redis *redis.Client
)

// GetRedisClient 获取redis客户端
func GetRedisClient() *redis.Client {
	return _redis
}

// SetRedisClient 设置redis客户端
func SetRedisClient(c *redis.Client) {
	if _redis != nil && _redis != c {
		_redis.Shutdown()
	}
	_redis = c
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

func (e RedisConnectOptions) GetRedisOptions() (*redis.Options, error) {
	r := &redis.Options{
		Network:    e.Network,
		Addr:       e.Addr,
		Username:   e.Username,
		Password:   e.Password,
		DB:         e.DB,
		MaxRetries: e.MaxRetries,
		PoolSize:   e.PoolSize,
	}
	var err error
	r.TLSConfig, err = getTLS(e.Tls)
	return r, err
}
