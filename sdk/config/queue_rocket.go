package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
)

const (
	RocketMqCfgName = "RocketMqConfig"
)

var insQueueRocket = cQueueRocket{}

type cQueueRocket struct {
	RocketOptions
	Urls []string
}

func QueueRocket() *cQueueRocket {
	return &insQueueRocket
}

func (c *cQueueRocket) String() string {
	return RocketMqCfgName
}

func (c *cQueueRocket) Init(ctx context.Context, s *Settings) error {
	_, err := c.GetRocketOptions()
	if err != nil {
		return err
	}
	c.Urls = []string{}
	return nil
}

// GetQueue get Rocket queue
func (c *cQueueRocket) GetQueue(ctx context.Context) (storage.AdapterQueue, error) {
	return queue.NewRocketMQ(
		ctx,
		c.Urls,
		&queue.RocketConsumerOptions{},
		&queue.RocketProducerOptions{},
		nil,
	)
}
