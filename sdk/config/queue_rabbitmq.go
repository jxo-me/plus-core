package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	queueLib "github.com/jxo-me/plus-core/core/queue"
	"github.com/jxo-me/plus-core/sdk/queue/rabbitmq"
)

const (
	RabbitmqQueueName = "rabbitmq"
)

var insQueueRabbit = cQueueRabbit{
	RabbitOptions: &RabbitOptions{},
}

type cQueueRabbit struct {
	*RabbitOptions
}

func QueueRabbit() *cQueueRabbit {
	return &insQueueRabbit
}

func (c *cQueueRabbit) String() string {
	return RabbitmqQueueName
}

func (c *cQueueRabbit) Init(ctx context.Context, s *Settings) error {
	var err error
	c.RabbitOptions, err = c.GetRabbitOptions(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

// GetQueue get Rabbit queue
func (c *cQueueRabbit) GetQueue(ctx context.Context) (queueLib.IQueue, error) {
	logger := glog.New()
	err := logger.SetConfigWithMap(g.Map{
		"flags":  glog.F_TIME_STD | glog.F_FILE_LONG,
		"path":   c.LogPath,
		"file":   c.LogFile,
		"level":  c.LogLevel,
		"stdout": c.LogStdout,
	})
	if err != nil {
		return nil, err
	}
	return rabbitmq.NewRabbitMQ(
		ctx,
		c.RabbitOptions.Dsn,
		c.RabbitOptions.ReconnectInterval,
		c.RabbitOptions.Cfg,
		logger,
	)
}
