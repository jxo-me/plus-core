package config

import (
	"context"
	"fmt"
	"github.com/jxo-me/rabbitmq-go"
	"time"
)

const (
	DefaultGroupName = "default" // Default configuration group name.
)

type RabbitOptions struct {
	Dsn      string `yaml:"dsn" json:"dsn"`
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	// Vhost specifies the namespace of permissions, exchanges, queues and
	// bindings on the server.  Dial sets this to the path parsed from the URL.
	Vhost             string        `yaml:"vhost" json:"vhost"`
	ChannelMax        int           `yaml:"channelMax" json:"channel_max"` // 0 max channels means 2^16 - 1
	FrameSize         int           `yaml:"frameSize" json:"frame_size"`   // 0 max bytes means unlimited
	Heartbeat         time.Duration `yaml:"heartbeat" json:"heartbeat"`    // less than 1s uses the server's interval
	ReconnectInterval int           `yaml:"reconnectInterval" json:"reconnectInterval"`
	// TLSClientConfig specifies the client configuration of the TLS connection
	// when establishing a tls transport.
	// If the URL uses an amqps scheme, then an empty tls.Config with the
	// ServerName from the URL is used.
	Tls *Tls             `json:"tls" yaml:"tls"`
	Cfg *rabbitmq.Config `json:"cfg" yaml:"cfg"`
}

func (e *RabbitOptions) GetReconnectInterval() int {
	return e.ReconnectInterval
}

func (e *RabbitOptions) GetDsn() string {
	if e.Dsn != "" {
		return e.Dsn
	}
	e.Dsn = fmt.Sprintf("amqp://%s:%s@%s/%s", e.Username, e.Password, e.Addr, e.Vhost)
	return e.Dsn
}

func (e *RabbitOptions) GetRabbitOptions(ctx context.Context, s *Settings) (*RabbitOptions, error) {
	// reconnectInterval
	e.Cfg = &rabbitmq.Config{}
	// cert key ca
	cert, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.cert", DefaultGroupName), "")
	if err != nil {
		return nil, err
	}
	key, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.key", DefaultGroupName), "")
	if err != nil {
		return nil, err
	}
	ca, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.ca", DefaultGroupName), "")
	if err != nil {
		return nil, err
	}
	tls := &Tls{
		Cert: cert.String(),
		Ca:   ca.String(),
		Key:  key.String(),
	}
	e.Cfg.TLSClientConfig, err = getTLS(tls)
	if err != nil {
		return nil, err
	}
	vhost, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.vhost", DefaultGroupName), "")
	if err != nil {
		return nil, err
	}
	if vhost.String() != "" {
		e.Cfg.Vhost = vhost.String()
	}
	channelMax, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.channelMax", DefaultGroupName), "0")
	if err != nil {
		return nil, err
	}
	reconnectInterval, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.reconnectInterval", DefaultGroupName), "5")
	if err != nil {
		return nil, err
	}
	e.ReconnectInterval = reconnectInterval.Int()

	if channelMax.Int() > 0 {
		e.Cfg.ChannelMax = channelMax.Int()
	}
	frameSize, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.frameSize", DefaultGroupName), "0")
	if err != nil {
		return nil, err
	}
	if frameSize.Int() > 0 {
		e.Cfg.FrameSize = frameSize.Int()
	}
	heartbeat, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.heartbeat", DefaultGroupName), "0")
	if err != nil {
		return nil, err
	}
	if heartbeat.Int() > 0 {
		e.Cfg.Heartbeat = time.Duration(heartbeat.Int()) * time.Second
	}
	dsn, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.dsn", DefaultGroupName), "")
	if err != nil {
		return nil, err
	}
	if dsn.String() != "" {
		e.Dsn = dsn.String()
	} else {
		addr, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.addr", DefaultGroupName), "")
		if err != nil {
			return nil, err
		}
		e.Addr = addr.String()
		username, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.username", DefaultGroupName), "")
		if err != nil {
			return nil, err
		}
		e.Username = username.String()
		password, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rabbitmq.%s.password", DefaultGroupName), "")
		if err != nil {
			return nil, err
		}
		e.Password = password.String()
		e.GetDsn()
	}

	return e, nil
}
