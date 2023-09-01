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
	// log
	LogPath   string `yaml:"logPath" json:"log_path"`
	LogFile   string `yaml:"logFile" json:"log_file"`
	LogLevel  string `yaml:"logLevel" json:"log_level"`
	LogStdout bool   `yaml:"logStdout" json:"log_stdout"`
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
