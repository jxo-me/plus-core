package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/glog"
	"io/ioutil"
)

var _redis *gredis.Redis

// GetRedisClient 获取redis客户端
func GetRedisClient() *gredis.Redis {
	return _redis
}

// SetRedisClient 设置redis客户端
func SetRedisClient(ctx context.Context, c *gredis.Redis) {
	if _redis != nil && _redis != c {
		err := _redis.Close(ctx)
		if err != nil {
			glog.Warning(ctx, "gRedis close error:", err)
		}
	}
	_redis = c
}

type RedisConnectOptions struct {
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	Tls      *Tls   `yaml:"tls" json:"tls"`
}

type Tls struct {
	Cert string `yaml:"cert" json:"cert"`
	Key  string `yaml:"key" json:"key"`
	Ca   string `yaml:"ca" json:"ca"`
}

func (e RedisConnectOptions) GetRedisOptions() (*gredis.Config, error) {
	r := &gredis.Config{
		Address: e.Addr,
		Pass:    e.Password,
		Db:      e.DB,
	}
	var err error
	r.TLSConfig, err = getTLS(e.Tls)
	return r, err
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
