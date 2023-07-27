package config

import (
	"context"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	nsq2 "github.com/jxo-me/plus-core/sdk/v2/queue/nsq"
	"github.com/nsqio/go-nsq"
)

const (
	NsqQueueName = "nsq"
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
	return NsqQueueName
}

func (c *cQueueNsq) Init(ctx context.Context) error {
	var err error
	c.Cfg, err = c.GetNsqOptions(ctx, Setting())
	if err != nil {
		return err
	}
	return nil
}

// GetQueue get NSQ queue
func (c *cQueueNsq) GetQueue(ctx context.Context) (queueLib.IQueue, error) {
	return nsq2.NewNSQ(c.Addresses, c.Cfg, c.ChannelPrefix)
}
