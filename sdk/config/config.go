package config

import (
	"fmt"
	"github.com/jxo-me/plus-core/pkg/v2/security"
	"github.com/jxo-me/rabbitmq-go"
	"time"
)

type Crypto struct {
	Enable    bool                      `json:"enable" yaml:"enable"`
	Algorithm string                    `json:"algorithm" yaml:"algorithm"`
	Rc4       security.Rc4CipherConfig  `json:"rc4" yaml:"rc4"`
	Rsa       security.RsaCiphersConfig `json:"rsa" yaml:"rsa"`
	Aes       security.AesCipherConfig  `json:"aes" yaml:"aes"`
}

type Metrics struct {
	Enable          bool      `json:"enable" yaml:"enable"`
	Path            string    `json:"path" yaml:"path"`
	SlowTime        int32     `json:"slowTime" yaml:"slowTime"`
	RequestDuration []float64 `json:"requestDuration" yaml:"requestDuration"`
}

type Jwt struct {
	Secret      string `yaml:"secret" json:"secret"`
	SigningKey  string `yaml:"signingKey" json:"signing_key"`
	Timeout     int64  `yaml:"timeout" json:"timeout"`
	MaxRefresh  int64  `yaml:"maxRefresh" json:"max_refresh"`
	IdentityKey string `yaml:"identityKey" json:"identity_key"`
}

type NSQOptions struct {
	DialTimeout  time.Duration `opt:"dial_timeout" default:"1s"`
	ReadTimeout  time.Duration `opt:"read_timeout" min:"100ms" max:"5m" default:"60s"`
	WriteTimeout time.Duration `opt:"write_timeout" min:"100ms" max:"5m" default:"1s"`

	// Addresses is the local address to use when dialing a nsqd.
	Addresses               []string      `opt:"addresses"`
	LookupdPollInterval     time.Duration `opt:"lookupd_poll_interval" min:"10ms" max:"5m" default:"60s"`
	LookupdPollJitter       float64       `opt:"lookupd_poll_jitter" min:"0" max:"1" default:"0.3"`
	MaxRequeueDelay         time.Duration `opt:"max_requeue_delay" min:"0" max:"60m" default:"15m"`
	DefaultRequeueDelay     time.Duration `opt:"default_requeue_delay" min:"0" max:"60m" default:"90s"`
	MaxBackoffDuration      time.Duration `opt:"max_backoff_duration" min:"0" max:"60m" default:"2m"`
	BackoffMultiplier       time.Duration `opt:"backoff_multiplier" min:"0" max:"60m" default:"1s"`
	MaxAttempts             uint16        `opt:"max_attempts" min:"0" max:"65535" default:"5"`
	LowRdyIdleTimeout       time.Duration `opt:"low_rdy_idle_timeout" min:"1s" max:"5m" default:"10s"`
	LowRdyTimeout           time.Duration `opt:"low_rdy_timeout" min:"1s" max:"5m" default:"30s"`
	RDYRedistributeInterval time.Duration `opt:"rdy_redistribute_interval" min:"1ms" max:"5s" default:"5s"`
	ClientID                string        `opt:"client_id"` // (defaults: short hostname)
	Hostname                string        `opt:"hostname"`
	UserAgent               string        `opt:"user_agent"`

	HeartbeatInterval time.Duration `opt:"heartbeat_interval" default:"30s"`
	SampleRate        int32         `opt:"sample_rate" min:"0" max:"99"`
	Tls               *Tls          `yaml:"tls" json:"tls"`

	// Compression Settings
	Deflate             bool          `opt:"deflate"`
	DeflateLevel        int           `opt:"deflate_level" min:"1" max:"9" default:"6"`
	Snappy              bool          `opt:"snappy"`
	OutputBufferSize    int64         `opt:"output_buffer_size" default:"16384"`
	OutputBufferTimeout time.Duration `opt:"output_buffer_timeout" default:"250ms"`
	MaxInFlight         int           `opt:"max_in_flight" min:"0" default:"1"`
	MsgTimeout          time.Duration `opt:"msg_timeout" min:"0"`
	AuthSecret          string        `opt:"auth_secret"`
}

type RabbitOptions struct {
	Dsn               string           `yaml:"dsn" json:"dsn"`
	Addr              string           `yaml:"addr" json:"addr"`
	Username          string           `yaml:"username" json:"username"`
	Password          string           `yaml:"password" json:"password"`
	Vhost             string           `yaml:"vhost" json:"vhost"`
	ChannelMax        int              `yaml:"channelMax" json:"channel_max"` // 0 max channels means 2^16 - 1
	FrameSize         int              `yaml:"frameSize" json:"frame_size"`   // 0 max bytes means unlimited
	Heartbeat         time.Duration    `yaml:"heartbeat" json:"heartbeat"`    // less than 1s uses the server's interval
	ReconnectInterval int              `yaml:"reconnectInterval" json:"reconnectInterval"`
	Tls               *Tls             `json:"tls" yaml:"tls"`
	Cfg               *rabbitmq.Config `json:"cfg" yaml:"cfg"`
	// log
	LogPath   string `yaml:"logPath" json:"log_path"`
	LogFile   string `yaml:"logFile" json:"log_file"`
	LogLevel  string `yaml:"logLevel" json:"log_level"`
	LogStdout bool   `yaml:"logStdout" json:"log_stdout"`
}

func (e *RabbitOptions) GetDsn() string {
	if e.Dsn != "" {
		return e.Dsn
	}
	e.Dsn = fmt.Sprintf("amqp://%s:%s@%s/%s", e.Username, e.Password, e.Addr, e.Vhost)
	return e.Dsn
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
