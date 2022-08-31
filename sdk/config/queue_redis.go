package config

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
	"github.com/robinjoseph08/redisqueue/v2"
	"time"
)

const (
	RedisQueueName = "redis"
)

var insQueueRedis = cQueueRedis{
	Producer: &redisqueue.ProducerOptions{},
	Consumer: &redisqueue.ConsumerOptions{},
}

type cQueueRedis struct {
	RedisConnectOptions
	Producer *redisqueue.ProducerOptions
	Consumer *redisqueue.ConsumerOptions
}

func QueueRedis() *cQueueRedis {
	return &insQueueRedis
}

func (c *cQueueRedis) String() string {
	return RedisQueueName
}

func (c *cQueueRedis) Init(ctx context.Context, s *Settings) error {
	client := Redis().GetClient()
	if client == nil {
		options, err := c.RedisConnectOptions.GetRedisOptions(ctx, s)
		if err != nil {
			return err
		}
		client = redis.NewClient(options)
		Redis().SetClient(ctx, client)
	}
	c.Producer.RedisClient = client
	c.Consumer.RedisClient = client
	c.Consumer.ReclaimInterval = c.Consumer.ReclaimInterval * time.Second
	c.Consumer.BlockingTimeout = c.Consumer.BlockingTimeout * time.Second
	c.Consumer.VisibilityTimeout = c.Consumer.VisibilityTimeout * time.Second
	return nil
}

// GetQueue get Redis queue
func (c *cQueueRedis) GetQueue(ctx context.Context) (storage.AdapterQueue, error) {
	return queue.NewRedis(c.Producer, c.Consumer)
}
