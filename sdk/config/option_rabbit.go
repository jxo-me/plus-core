package config

import (
	"fmt"
	"github.com/jxo-me/rabbitmq-go"
	"time"
)

type RabbitOptions struct {
	Dsn      string `yaml:"dsn" json:"dsn"`
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	// Vhost specifies the namespace of permissions, exchanges, queues and
	// bindings on the server.  Dial sets this to the path parsed from the URL.
	Vhost      string        `yaml:"vhost" json:"vhost"`
	ChannelMax int           `yaml:"channel_max" json:"channel_max"` // 0 max channels means 2^16 - 1
	FrameSize  int           `yaml:"frame_size" json:"frame_size"`   // 0 max bytes means unlimited
	Heartbeat  time.Duration `yaml:"heartbeat" json:"heartbeat"`     // less than 1s uses the server's interval
	// TLSClientConfig specifies the client configuration of the TLS connection
	// when establishing a tls transport.
	// If the URL uses an amqps scheme, then an empty tls.Config with the
	// ServerName from the URL is used.
	Tls *Tls
}

func (e RabbitOptions) GetDsn() string {
	e.Dsn = fmt.Sprintf("amqp://%s:%s@%s/%s", e.Username, e.Password, e.Addr, e.Vhost)
	return e.Dsn
}

func (e RabbitOptions) GetRabbitOptions() (*rabbitmq.Config, error) {
	cfg := &rabbitmq.Config{}
	var err error
	cfg.TLSClientConfig, err = getTLS(e.Tls)
	if err != nil {
		return nil, err
	}
	if e.Vhost != "" {
		cfg.Vhost = e.Vhost
	}
	if e.ChannelMax > 0 {
		cfg.ChannelMax = e.ChannelMax
	}
	if e.FrameSize > 0 {
		cfg.FrameSize = e.FrameSize
	}
	if e.Heartbeat > 0 {
		cfg.Heartbeat = e.Heartbeat * time.Second
	}

	return cfg, nil
}
