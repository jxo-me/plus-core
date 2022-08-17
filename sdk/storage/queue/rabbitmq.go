package queue

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/rabbitmq-go"
)

func NewRabbitMQ(
	ctx context.Context,
	dsn string,
	routingKeys []string,
	exchange string,
	exchangeType string,
	cfg *rabbitmq.Config,
	consumerOptions *rabbitmq.ConsumerOptions,
	publisherOptions *rabbitmq.PublisherOptions,
) (*RabbitMQ, error) {
	var err error
	var consumer rabbitmq.Consumer
	r := &RabbitMQ{
		Url:              dsn,
		RoutingKeys:      routingKeys,
		Exchange:         exchange,
		ExchangeType:     exchangeType,
		ConsumerOptions:  consumerOptions,
		PublisherOptions: publisherOptions,
	}
	if cfg != nil {
		r.Config = *cfg
	}
	consumer, err = r.newConsumer(ctx, r.ConsumerOptions)
	if err != nil {
		return nil, err
	}
	r.consumer = &consumer
	r.producer, err = r.newProducer(ctx, r.PublisherOptions)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// RabbitMQ cache implement
type RabbitMQ struct {
	Url              string
	Handler          []rabbitmq.Handler
	RoutingKeys      []string
	Exchange         string
	ExchangeType     string
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
	// exchange exchangeType routingKey
	exchange := "events"
	if message.GetGroupId() != "" {
		exchange = message.GetGroupId()
	}
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	err = r.producer.Publish(
		ctx,
		rb,
		[]string{message.GetRoutingKey()},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(exchange),
	)
	return err
}

// Consumer 监听消费者
func (r *RabbitMQ) Consumer(ctx context.Context, queueName string, f storage.ConsumerFunc) {
	// exchange exchangeType routingKey
	err := r.consumer.StartConsuming(ctx,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			glog.Debug(ctx, "rabbitmq consumed: %s\n", string(d.Body))
			m := new(Message)
			m.SetValues(gconv.Map(d.Body))
			m.SetRoutingKey(d.RoutingKey)
			m.SetId(d.MessageId)
			err := f(ctx, m)
			if err != nil {
				return rabbitmq.NackRequeue
			}
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		queueName,
		r.RoutingKeys,
		rabbitmq.WithConsumeOptionsConsumerName(queueName),
		rabbitmq.WithConsumeOptionsBindingExchangeName(r.Exchange),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(r.ExchangeType),
		rabbitmq.WithConsumeOptionsConcurrency(1), // goroutine num
		rabbitmq.WithConsumeOptionsConsumerAutoAck(true),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsQueueDurable,
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
