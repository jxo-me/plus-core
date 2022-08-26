package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
	"github.com/nsqio/go-nsq"
)

const (
	NSQMqCfgName = "nsqMqConfig"
)

var insQueueNSQ = cQueueNsq{
	Cfg: &nsq.Config{},
}

type cQueueNsq struct {
	Cfg *nsq.Config
	NSQOptions
	ChannelPrefix string
}

func QueueNsq() *cQueueNsq {
	return &insQueueNSQ
}

func (c *cQueueNsq) String() string {
	return NSQMqCfgName
}

func (c *cQueueNsq) Init(ctx context.Context, s *Settings) error {
	var err error
	c.Cfg, err = c.GetNsqOptions(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

// GetQueue get NSQ queue
func (c *cQueueNsq) GetQueue(ctx context.Context) (storage.AdapterQueue, error) {
	return queue.NewNSQ(c.Addresses, c.Cfg, c.ChannelPrefix)
}
