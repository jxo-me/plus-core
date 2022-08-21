package queue

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/robinjoseph08/redisqueue/v2"

	"github.com/jxo-me/plus-core/sdk/storage"
)

// NewRedis redis模式
func NewRedis(producerOptions *redisqueue.ProducerOptions,
	consumerOptions *redisqueue.ConsumerOptions) (*Redis, error) {
	var err error
	r := &Redis{}
	r.producer, err = r.newProducer(producerOptions)
	if err != nil {
		return nil, err
	}
	r.consumer, err = r.newConsumer(consumerOptions)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis cache implement
type Redis struct {
	client   *redis.Client
	consumer *redisqueue.Consumer
	producer *redisqueue.Producer
}

func (Redis) String() string {
	return "redis"
}

func (r *Redis) newConsumer(options *redisqueue.ConsumerOptions) (*redisqueue.Consumer, error) {
	if options == nil {
		options = &redisqueue.ConsumerOptions{}
	}
	return redisqueue.NewConsumerWithOptions(options)
}

func (r *Redis) newProducer(options *redisqueue.ProducerOptions) (*redisqueue.Producer, error) {
	if options == nil {
		options = &redisqueue.ProducerOptions{}
	}
	return redisqueue.NewProducerWithOptions(options)
}

// Publish 消息入生产者
func (r *Redis) Publish(ctx context.Context, message storage.Messager, optionFuncs ...func(*storage.PublishOptions)) error {
	err := r.producer.Enqueue(&redisqueue.Message{
		ID:     message.GetId(),
		Stream: message.GetRoutingKey(),
		Values: message.GetValues(),
	})
	return err
}

// Consumer 监听消费者
func (r *Redis) Consumer(ctx context.Context, name string, f storage.ConsumerFunc, optionFuncs ...func(*storage.ConsumeOptions)) {
	r.consumer.Register(name, func(message *redisqueue.Message) error {
		m := new(Message)
		m.SetValues(message.Values)
		m.SetRoutingKey(message.Stream)
		m.SetId(message.ID)
		return f(ctx, m)
	})
}

func (r *Redis) Run(ctx context.Context) {
	r.consumer.Run()
}

func (r *Redis) Shutdown(ctx context.Context) {
	r.consumer.Shutdown()
}
