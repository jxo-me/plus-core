package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/queue/rabbitmq"
	rabbitmqGo "github.com/jxo-me/rabbitmq-go"
)

const (
	RabbitmqQueueName = "rabbitmq"
)

var insQueueRabbit = cQueueRabbit{
	List: map[string]*RabbitOptions{},
}

type cQueueRabbit struct {
	List map[string]*RabbitOptions
}

func QueueRabbit() *cQueueRabbit {
	return &insQueueRabbit
}

func (c *cQueueRabbit) String() string {
	return RabbitmqQueueName
}

func (c *cQueueRabbit) Init(ctx context.Context) error {
	var err error
	conf, err := Setting().Cfg().Get(ctx, "settings.queue.rabbitmq", "")
	if err != nil {
		return err
	}
	list := make(map[string]*RabbitOptions)
	err = conf.Scan(&list)
	if err != nil {
		return err
	}
	for key, config := range list {
		if config.Tls != nil {
			tls := &Tls{
				Cert: config.Tls.Cert,
				Ca:   config.Tls.Ca,
				Key:  config.Tls.Key,
			}
			if config.Cfg == nil {
				config.Cfg = &rabbitmqGo.Config{}
			}
			config.Cfg.TLSClientConfig, err = getTLS(tls)
			if err != nil {
				return err
			}
		}
		if config.Dsn == "" {
			config.GetDsn()
		}
		c.List[key] = config
	}

	return nil
}

// GetQueue get Rabbit queue
func (c *cQueueRabbit) GetQueue(ctx context.Context) (map[string]queueLib.IQueue, error) {
	list := make(map[string]queueLib.IQueue)
	for key, options := range c.List {
		logger := glog.New()
		err := logger.SetConfigWithMap(g.Map{
			"flags":  glog.F_TIME_STD | glog.F_FILE_LONG,
			"path":   options.LogPath,
			"file":   options.LogFile,
			"level":  options.LogLevel,
			"stdout": options.LogStdout,
		})
		if err != nil {
			return nil, err
		}
		list[key], err = rabbitmq.NewRabbitMQ(
			ctx,
			options.Dsn,
			options.ReconnectInterval,
			options.Cfg,
			logger,
		)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
