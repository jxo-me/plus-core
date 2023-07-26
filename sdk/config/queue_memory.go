package config

import (
	"context"
	queueLib "github.com/jxo-me/plus-core/core/queue"
	"github.com/jxo-me/plus-core/sdk/queue/memory"
)

const (
	MemoryQueueName = "memory"
)

var insQueueMemory = cQueueMemory{}

type cQueueMemory struct {
	PoolSize uint `json:"poolSize" yaml:"poolSize"`
}

func QueueMemory() *cQueueMemory {
	return &insQueueMemory
}

func (c *cQueueMemory) String() string {
	return MemoryQueueName
}

func (c *cQueueMemory) Init(ctx context.Context, s *Settings) error {
	poolSize, err := s.Cfg().Get(ctx, "settings.queue.memory.poolSize", 10000)
	if err != nil {
		return err
	}
	c.PoolSize = poolSize.Uint()
	return nil
}

// GetQueue get Memory queue
func (c *cQueueMemory) GetQueue(ctx context.Context) (queueLib.IQueue, error) {
	return memory.NewMemory(c.PoolSize), nil
}
