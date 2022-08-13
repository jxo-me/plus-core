package config

// Amqp
// Config is the config of amqp client and has some custom options
type Amqp struct {
	DSN string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`

	MaxChannelsPerConnection int `mapstructure:"maxChannelsPerConnection" json:"maxChannelsPerConnection" yaml:"maxChannelsPerConnection"`
	MaxIdleChannels          int `mapstructure:"maxIdleChannels" json:"maxIdleChannels" yaml:"maxIdleChannels"`
	MaxConnections           int `mapstructure:"maxConnections" json:"maxConnections" yaml:"maxConnections"`
	MinConnections           int `mapstructure:"minConnections" json:"minConnections" yaml:"minConnections"`

	Debug bool `mapstructure:"debug" json:"debug" yaml:"debug"`
}
