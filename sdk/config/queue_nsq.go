package config

import (
	"context"
	"github.com/gogf/gf/v2/os/gcfg"
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

func (c *cQueueNsq) Init(ctx context.Context, s *gcfg.Config) error {
	var err error
	c.Cfg, err = c.GetNsqOptions(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

// GetQueue get NSQ queue
func (c *cQueueNsq) GetQueue(ctx context.Context) (map[string]queueLib.IQueue, error) {
	var err error
	list := make(map[string]queueLib.IQueue)
	list[DefaultGroupName], err = nsq2.NewNSQ(c.Addresses, c.Cfg, c.ChannelPrefix)
	return list, err
}
