package config

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
	"github.com/jxo-me/rabbitmq-go"
	"github.com/robinjoseph08/redisqueue/v2"
	"time"
)

var insQueue = Queue{}

type Queue struct {
	Redis  *QueueRedis  `json:"redis" yaml:"redis"`
	Memory *QueueMemory `json:"memory" yaml:"memory"`
	Rabbit *QueueRabbit `json:"rabbit" yaml:"rabbit"`
	Rocket *QueueRocket `json:"rocket" yaml:"rocket"`
	NSQ    *QueueNSQ    `json:"nsq" yaml:"nsq"`
}

type QueueRedis struct {
	RedisConnectOptions
	Producer *redisqueue.ProducerOptions
	Consumer *redisqueue.ConsumerOptions
}

type QueueMemory struct {
	PoolSize uint
}

type QueueRabbit struct {
	RabbitOptions
	Producer *rabbitmq.PublisherOptions
	Consumer *rabbitmq.ConsumeOptions
}

type QueueRocket struct {
	RocketOptions
}

type QueueNSQ struct {
	NSQOptions
	ChannelPrefix string
}

func QueueConfig() *Queue {
	return &insQueue
}

// Empty 空设置
func (e *Queue) Empty() bool {
	return e.Memory == nil && e.Redis == nil && e.NSQ == nil && e.Rabbit == nil && e.Rocket == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e *Queue) Setup(ctx context.Context, s *Settings) (storage.AdapterQueue, error) {
	if e.Redis != nil {
		e.Redis.Consumer.ReclaimInterval = e.Redis.Consumer.ReclaimInterval * time.Second
		e.Redis.Consumer.BlockingTimeout = e.Redis.Consumer.BlockingTimeout * time.Second
		e.Redis.Consumer.VisibilityTimeout = e.Redis.Consumer.VisibilityTimeout * time.Second
		client := GetRedisClient()
		if client == nil {
			options, err := e.Redis.RedisConnectOptions.GetRedisOptions()
			if err != nil {
				return nil, err
			}
			client = redis.NewClient(options)
			_redis = client
		}
		e.Redis.Producer.RedisClient = client
		e.Redis.Consumer.RedisClient = client
		return queue.NewRedis(e.Redis.Producer, e.Redis.Consumer)
	}
	// rabbitmq queue
	if e.Rabbit != nil {
		_, err := e.Rabbit.GetRabbitOptions()
		if err != nil {
			return nil, err
		}
		dsn := e.Rabbit.GetDsn()
		return queue.NewRabbitMQ(
			ctx, dsn, &rabbitmq.Config{},
		)
	}
	// rocketmq queue
	if e.Rocket != nil {
		_, err := e.Rocket.GetRocketOptions()
		if err != nil {
			return nil, err
		}
		urls := []string{}
		return queue.NewRocketMQ(
			ctx,
			urls,
			&queue.RocketConsumerOptions{},
			&queue.RocketProducerOptions{},
			nil,
		)
	}
	// NSQ
	if e.NSQ != nil {
		cfg, err := e.NSQ.GetNSQOptions()
		if err != nil {
			return nil, err
		}
		return queue.NewNSQ(e.NSQ.Addresses, cfg, e.NSQ.ChannelPrefix)
	}
	return queue.NewMemory(e.Memory.PoolSize), nil
}
