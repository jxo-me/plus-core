package queue

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/rabbitmq-go"
	"log"
)

// RabbitMQ cache implement
type RabbitMQ struct {
	Url              string
	Config           rabbitmq.Config
	consumer         *rabbitmq.Consumer
	ConsumerOptions  *rabbitmq.ConsumerOptions
	producer         *rabbitmq.Publisher
	PublisherOptions *rabbitmq.PublisherOptions
}

func (RabbitMQ) String() string {
	return "rabbitmq"
}

func (r *RabbitMQ) newConsumer(ctx context.Context, options *rabbitmq.ConsumerOptions) (rabbitmq.Consumer, error) {
	if options == nil {
		r.ConsumerOptions = &rabbitmq.ConsumerOptions{}
	}
	return rabbitmq.NewConsumer(ctx, r.Url, r.Config,
		rabbitmq.WithConsumerOptionsLogger(r.ConsumerOptions.Logger),
		rabbitmq.WithConsumerOptionsReconnectInterval(r.ConsumerOptions.ReconnectInterval),
	)
}

func (r *RabbitMQ) newProducer(ctx context.Context, options *rabbitmq.PublisherOptions) (*rabbitmq.Publisher, error) {
	if options == nil {
		r.PublisherOptions = &rabbitmq.PublisherOptions{}
	}
	return rabbitmq.NewPublisher(ctx, r.Url, r.Config,
		rabbitmq.WithPublisherOptionsLogger(r.PublisherOptions.Logger),
		rabbitmq.WithPublisherOptionsReconnectInterval(r.PublisherOptions.ReconnectInterval),
	)
}

// Publish 消息入生产者
func (r *RabbitMQ) Publish(ctx context.Context, message storage.Messager) error {
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	err = r.producer.Publish(
		ctx,
		rb,
		[]string{message.GetStream()},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	return err
}

// Consumer 监听消费者
func (r *RabbitMQ) Consumer(ctx context.Context, name string, f storage.ConsumerFunc) {
	err := r.consumer.StartConsuming(ctx,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			m := new(Message)
			m.SetValues(gconv.Map(d.Body))
			m.SetStream(d.RoutingKey)
			m.SetID(d.MessageId)
			err := f(ctx, m)
			if err != nil {
				return rabbitmq.NackRequeue
			}
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		name,
		[]string{"routing_key", "routing_key_2"},
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName("events"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConsumerName(name),
	)
	if err != nil {
		glog.Errorf(ctx, "rabbitmq consumer StartConsuming error:%v", err)
		return
	}
}

func (r *RabbitMQ) Run(ctx context.Context) {
	return
}

func (r *RabbitMQ) Shutdown(ctx context.Context) {
	if r.producer != nil {
		err := r.producer.Close(ctx)
		if err != nil {
			glog.Warning(ctx, "rabbitmq producer Close error", err)
		}
	}
	if r.consumer != nil {
		err := r.consumer.Close(ctx)
		if err != nil {
			glog.Warning(ctx, "rabbitmq consumer Close error", err)
		}
	}
}
