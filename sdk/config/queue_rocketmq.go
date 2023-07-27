package config

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/queue/rocketmq"
)

const (
	RocketQueueName = "rocketmq"
)

var insQueueRocket = cQueueRocket{
	RocketOptions: &RocketOptions{},
}

type cQueueRocket struct {
	*RocketOptions
}

func QueueRocket() *cQueueRocket {
	return &insQueueRocket
}

func (c *cQueueRocket) String() string {
	return RocketQueueName
}

func (c *cQueueRocket) Init(ctx context.Context) error {
	var err error
	c.RocketOptions, err = c.GetRocketOptions(ctx, Setting())
	if err != nil {
		return err
	}
	// primitive.Credentials
	if c.AccessKey != "" && c.SecretKey != "" {
		c.Credentials = &primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}
	}
	return nil
}

// GetQueue get Rocket queue
func (c *cQueueRocket) GetQueue(ctx context.Context) (queueLib.IQueue, error) {
	logger := glog.New()
	err := logger.SetConfigWithMap(g.Map{
		"path":   c.LogPath,
		"file":   c.LogFile,
		"level":  c.LogLevel,
		"stdout": c.LogStdout,
	})
	if err != nil {
		return nil, err
	}
	return rocketmq.NewRocketMQ(
		ctx,
		c.Urls,
		c.Credentials,
		logger,
	)
}
