package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
)

const (
	RabbitmqCfgName = "RabbitMqConfig"
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
	return RabbitmqCfgName
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
func (c *cQueueRabbit) GetQueue(ctx context.Context) (storage.AdapterQueue, error) {
	return queue.NewRabbitMQ(
		ctx, c.RabbitOptions.Dsn, c.RabbitOptions.ReconnectInterval, c.RabbitOptions.Cfg,
	)
}
