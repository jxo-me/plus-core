package config

import "context"

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// Init the config
	Init(ctx context.Context, opts ...Option) error
	// Options in the config
	Options() Options
}

// Options 配置的参数
type Options struct {
	// for alternative data
	Context context.Context
}

// Option 调用类型
type Option func(o *Options)
