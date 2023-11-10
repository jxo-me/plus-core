package config

import (
	"context"
	"github.com/gogf/gf/v2/os/gcfg"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/queue/memory"
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

func (c *cQueueMemory) Init(ctx context.Context, s *gcfg.Config) error {
	poolSize, err := s.Get(ctx, "settings.queue.memory.poolSize", 10000)
	if err != nil {
		return err
	}
	c.PoolSize = poolSize.Uint()
	return nil
}

// GetQueue get Memory queue
func (c *cQueueMemory) GetQueue(ctx context.Context) (map[string]queueLib.IQueue, error) {
	list := make(map[string]queueLib.IQueue)
	list[DefaultGroupName] = memory.NewMemory(c.PoolSize)
	return list, nil
}
