package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/glog"
	"io/ioutil"
)

var (
	_gredis *gredis.Redis
	_redis  *redis.Client
)

type Tls struct {
	Cert string `yaml:"cert" json:"cert"`
	Key  string `yaml:"key" json:"key"`
	Ca   string `yaml:"ca" json:"ca"`
}

func getTLS(c *Tls) (*tls.Config, error) {
	if c != nil && c.Cert != "" {
		// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
		cert, err := tls.LoadX509KeyPair(c.Cert, c.Key)
		if err != nil {
			fmt.Printf("tls.LoadX509KeyPair err: %v\n", err)
			return nil, err
		}
		// 创建一个新的、空的 CertPool，并尝试解析 PEM 编码的证书，解析成功会将其加到 CertPool 中
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(c.Ca)
		if err != nil {
			fmt.Printf("ioutil.ReadFile err: %v\n", err)
			return nil, err
		}

		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			fmt.Println("certPool.AppendCertsFromPEM err")
			return nil, err
		}
		return &tls.Config{
			// 设置证书链，允许包含一个或多个
			Certificates: []tls.Certificate{cert},
			// 要求必须校验客户端的证书
			ClientAuth: tls.RequireAndVerifyClientCert,
			// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
			ClientCAs: certPool,
		}, nil
	}
	return nil, nil
}

// GetGRedisClient 获取redis客户端
func GetGredisClient() *gredis.Redis {
	return _gredis
}

// SetGredisClient 设置redis客户端
func SetGredisClient(ctx context.Context, c *gredis.Redis) {
	if _gredis != nil && _gredis != c {
		err := _gredis.Close(ctx)
		if err != nil {
			glog.Warning(ctx, "gRedis close error:", err)
		}
	}
	_gredis = c
}

type GredisConnectOptions struct {
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	Tls      *Tls   `yaml:"tls" json:"tls"`
}

func (e GredisConnectOptions) GetGredisOptions() (*gredis.Config, error) {
	r := &gredis.Config{
		Address: e.Addr,
		Pass:    e.Password,
		Db:      e.DB,
	}
	var err error
	r.TLSConfig, err = getTLS(e.Tls)
	return r, err
}

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
