package config

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/net/ghttp"
	service "github.com/jxo-me/plus-core/pkg/v2/erlang"
	"github.com/jxo-me/plus-core/pkg/v2/security"
	"github.com/jxo-me/plus-core/pkg/v2/ws"
	"github.com/jxo-me/rabbitmq-go"
	"time"
)

type Config struct {
	Database map[string]*gdb.ConfigNode     `json:"database"`
	Redis    map[string]*RedisOptions       `json:"redis"`
	Server   *ghttp.ServerConfig            `json:"server"`
	Settings *SettingOptions                `json:"settings"`
	Bot      *BotGroups                     `json:"bot,omitempty"`
	Erlang   map[string]*service.NodeConfig `json:"erlang,omitempty"`
}

type RedisOptions struct {
	gredis.Config
	// Base number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	// If there are not enough connections in the pool, new connections will be allocated in excess of PoolSize;
	// you can limit it through MaxActiveConns
	PoolSize int
}

type CryptoOptions struct {
	Enable    bool                      `json:"enable" yaml:"enable"`
	Algorithm string                    `json:"algorithm" yaml:"algorithm"`
	Rc4       security.Rc4CipherConfig  `json:"rc4" yaml:"rc4"`
	Rsa       security.RsaCiphersConfig `json:"rsa" yaml:"rsa"`
	Aes       security.AesCipherConfig  `json:"aes" yaml:"aes"`
}

type MetricsOptions struct {
	Enable          bool      `json:"enable" yaml:"enable"`
	Path            string    `json:"path" yaml:"path"`
	SlowTime        int32     `json:"slowTime" yaml:"slowTime"`
	RequestDuration []float64 `json:"requestDuration" yaml:"requestDuration"`
}

type JwtOptions struct {
	Secret          string `yaml:"secret" json:"secret"`
	SigningKey      string `yaml:"signingKey" json:"signing_key"`
	Timeout         int64  `yaml:"timeout" json:"timeout"`
	MaxRefresh      int64  `yaml:"maxRefresh" json:"max_refresh"`
	IdentityKey     string `yaml:"identityKey" json:"identity_key"`
	BlacklistPrefix string `yaml:"blacklistPrefix" json:"blacklistPrefix"`
}

type BotGroups struct {
	Default *TgBot `json:"default,omitempty"`
}

type TgBot struct {
	LogPath   string `json:"logPath,omitempty"`
	LogFile   string `json:"logFile,omitempty"`
	LogLevel  string `json:"logLevel,omitempty"`
	LogStdout bool   `json:"logStdout,omitempty"`
	HookUrl   string `json:"hook,omitempty"`
	Token     string `json:"token,omitempty"`
}

type SettingOptions struct {
	System      *System                        `json:"system,omitempty"`
	Queue       *QueueGroups                   `json:"queue,omitempty"`
	Uploads     *UploadGroups                  `json:"uploads,omitempty"`
	Auth        map[string]*JwtAuth            `json:"auth,omitempty"`
	Websocket   *ws.Config                     `json:"websocket,omitempty"`
	FailedLimit map[string]*FailedLimitOptions `json:"failedLimit,omitempty"`
	Metrics     *MetricsOptions                `json:"metrics,omitempty"`
}

type JwtAuth struct {
	Jwt *JwtOptions `json:"jwt,omitempty"`
}

type FailedLimitOptions struct {
	Limit  int64 `json:"limit,omitempty"`
	Expire int64 `json:"expire,omitempty"`
}

type QueueGroups struct {
	Memory   *MemoryOptions              `json:"memory,omitempty"`
	Rabbitmq map[string]*RabbitmqOptions `json:"rabbitmq,omitempty"`
	Rocketmq map[string]*RocketmqOptions `json:"rocketmq,omitempty"`
	Nsq      *NSQOptions                 `json:"nsq"`
}

type MemoryOptions struct {
	PoolSize uint `json:"poolSize,omitempty"`
}

type RabbitmqOptions struct {
	DSN               string           `json:"dsn,omitempty"`
	Vhost             string           `json:"vhost,omitempty"`
	ChannelMax        int64            `json:"channelMax,omitempty"`
	FrameSize         int64            `json:"frameSize,omitempty"`
	Heartbeat         int64            `json:"heartbeat,omitempty"`
	ReconnectInterval int              `json:"reconnectInterval,omitempty"`
	Cfg               *rabbitmq.Config `json:"cfg" yaml:"cfg"`
	LogPath           string           `json:"logPath,omitempty"`
	LogFile           string           `json:"logFile,omitempty"`
	LogLevel          string           `json:"logLevel,omitempty"`
	LogStdout         bool             `json:"logStdout,omitempty"`
}

type RocketmqOptions struct {
	Urls []string `json:"urls,omitempty"`
	*primitive.Credentials
	LogPath   string `json:"logPath,omitempty"`
	LogFile   string `json:"logFile,omitempty"`
	LogLevel  string `json:"logLevel,omitempty"`
	LogStdout bool   `json:"logStdout,omitempty"`
}

type NSQOptions struct {
	DialTimeout time.Duration `opt:"dial_timeout" default:"1s"`

	// Deadlines for network reads and writes
	ReadTimeout  time.Duration `opt:"read_timeout" min:"100ms" max:"5m" default:"60s"`
	WriteTimeout time.Duration `opt:"write_timeout" min:"100ms" max:"5m" default:"1s"`
	// Addresses is the local address to use when dialing an nsqd.
	Addresses     []string `opt:"addresses"`
	ChannelPrefix string
}

type System struct {
	EnableCAPTCHA  bool           `json:"enableCaptcha,omitempty"`
	Logger         *SystemLogger  `json:"logger,omitempty"`
	DataPermission bool           `json:"dataPermission,omitempty"`
	Crypto         *CryptoOptions `json:"crypto,omitempty"`
}

type SystemLogger struct {
	EnabledDB bool `json:"enabledDb,omitempty"`
}

type UploadGroups struct {
	Normal *UploadOptions         `json:"normal,omitempty"`
	Tus    map[string]*TusOptions `json:"tus,omitempty"`
}

type UploadOptions struct {
	Path     string `json:"path,omitempty"`
	BasePath string `json:"basePath,omitempty"`
}

type TusOptions struct {
	Path      string `json:"path,omitempty"`
	BasePath  string `json:"basePath,omitempty"`
	LogPath   string `json:"logPath,omitempty"`
	LogFile   string `json:"logFile,omitempty"`
	LogLevel  string `json:"logLevel,omitempty"`
	LogStdout bool   `json:"logStdout,omitempty"`
}
