package parsing

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/config"
	"github.com/jxo-me/plus-core/sdk/v2/queue/memory"
	nsq2 "github.com/jxo-me/plus-core/sdk/v2/queue/nsq"
	"github.com/jxo-me/plus-core/sdk/v2/queue/rabbitmq"
	"github.com/jxo-me/plus-core/sdk/v2/queue/rocketmq"
	"github.com/nsqio/go-nsq"
)

func ParseRabbitMQ(ctx context.Context, cfg *config.RabbitmqOptions) (queueLib.IQueue, error) {
	logger := glog.New()
	err := logger.SetConfigWithMap(g.Map{
		"flags":  glog.F_TIME_STD | glog.F_FILE_LONG,
		"path":   cfg.LogPath,
		"file":   cfg.LogFile,
		"level":  cfg.LogLevel,
		"stdout": cfg.LogStdout,
	})

	q, err := rabbitmq.NewRabbitMQ(
		ctx,
		cfg.DSN,
		cfg.MaxConnections,
		cfg.ReconnectInterval,
		cfg.Cfg,
		logger,
	)

	return q, err
}

func ParseRocketMQ(ctx context.Context, cfg *config.RocketmqOptions) (queueLib.IQueue, error) {
	logger := glog.New()
	err := logger.SetConfigWithMap(g.Map{
		"flags":  glog.F_TIME_STD | glog.F_FILE_LONG,
		"path":   cfg.LogPath,
		"file":   cfg.LogFile,
		"level":  cfg.LogLevel,
		"stdout": cfg.LogStdout,
	})
	if err != nil {
		return nil, err
	}

	// primitive.Credentials
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		cfg.Credentials = &primitive.Credentials{
			AccessKey: cfg.AccessKey,
			SecretKey: cfg.SecretKey,
		}
	}
	q, err := rocketmq.NewRocketMQ(
		ctx,
		cfg.Urls,
		cfg.Credentials,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func ParseNSQMQ(ctx context.Context, cfg *config.NSQOptions) (queueLib.IQueue, error) {
	c := &nsq.Config{}
	q, err := nsq2.NewNSQ(cfg.Addresses, c, cfg.ChannelPrefix)
	return q, err
}

func ParseMemoryMQ(ctx context.Context, cfg *config.MemoryOptions) (queueLib.IQueue, error) {
	return memory.NewMemory(cfg.PoolSize), nil
}
