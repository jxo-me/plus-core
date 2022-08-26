package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
)

const (
	MemoryMqCfgName = "memoryMqConfig"
)

var insQueueMemory = cQueueMemory{}

type cQueueMemory struct {
	PoolSize uint `json:"poolSize" yaml:"poolSize"`
}

func QueueMemory() *cQueueMemory {
	return &insQueueMemory
}

func (c *cQueueMemory) String() string {
	return MemoryMqCfgName
}

func (c *cQueueMemory) Init(ctx context.Context, s *Settings) error {
	poolSize, err := s.Cfg().Get(ctx, "memoryQueue.poolSize", 10000)
	if err != nil {
		return err
	}
	c.PoolSize = poolSize.Uint()
	return nil
}

// GetQueue get Memory queue
func (c *cQueueMemory) GetQueue(ctx context.Context) (storage.AdapterQueue, error) {
	return queue.NewMemory(c.PoolSize), nil
}
