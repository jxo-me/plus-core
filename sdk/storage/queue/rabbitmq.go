package queue

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/rabbitmq-go"
	"time"
)

func NewRabbitMQ(
	ctx context.Context,
	dsn string,
	reconnectInterval int,
	cfg *rabbitmq.Config,
) (*RabbitMQ, error) {
	//var err error
	//var consumer rabbitmq.Consumer
	r := &RabbitMQ{
		Url:               dsn,
		ReconnectInterval: reconnectInterval,
	}
	if cfg != nil {
		r.Config = *cfg
	}
	return r, nil
}

// RabbitMQ cache implement
type RabbitMQ struct {
	Url               string
	ReconnectInterval int
	Handler           []rabbitmq.Handler
	Config            rabbitmq.Config
	consumer          *rabbitmq.Consumer
	ConsumerOptions   *rabbitmq.ConsumerOptions
	producer          *rabbitmq.Publisher
	PublisherOptions  *rabbitmq.PublisherOptions
}

func (RabbitMQ) String() string {
	return "rabbitmq"
}

func (r *RabbitMQ) newConsumer(ctx context.Context) (rabbitmq.Consumer, error) {
	return rabbitmq.NewConsumer(ctx,
		r.Url,
		r.Config,
		rabbitmq.WithConsumerOptionsLogger(g.Log()),
		rabbitmq.WithConsumerOptionsReconnectInterval(time.Duration(r.ReconnectInterval)*time.Second),
	)
}

func (r *RabbitMQ) newProducer(ctx context.Context) (*rabbitmq.Publisher, error) {
	return rabbitmq.NewPublisher(ctx,
		r.Url,
		r.Config,
		rabbitmq.WithPublisherOptionsLogger(g.Log()),
		rabbitmq.WithPublisherOptionsReconnectInterval(time.Duration(r.ReconnectInterval)*time.Second),
	)
}

// Publish 消息入生产者
func (r *RabbitMQ) Publish(ctx context.Context, message storage.Messager, optionFuncs ...func(*storage.PublishOptions)) error {
	// exchange exchangeType routingKey
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	options := &storage.PublishOptions{
		ContentType: "application/json",
		MessageID:   gctx.CtxId(ctx),
	}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if r.producer == nil {
		// connection ..
		r.producer, err = r.newProducer(ctx)
		if err != nil {
			glog.Warning(ctx, "rabbitmq newProducer error:", err)
			return err
		}
	}

	err = r.producer.Publish(
		ctx,
		rb,
		[]string{message.GetRoutingKey()},
		rabbitmq.WithPublishOptionsExchange(options.Exchange),
		rabbitmq.WithPublishOptionsContentType(options.ContentType),
		rabbitmq.WithPublishOptionsMessageID(options.MessageID),
		rabbitmq.WithPublishOptionsAppID(options.AppID),
		rabbitmq.WithPublishOptionsUserID(options.UserID),
		rabbitmq.WithPublishOptionsReplyTo(options.ReplyTo),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
	return err
}

// Consumer 监听消费者
func (r *RabbitMQ) Consumer(ctx context.Context, queueName string, f storage.ConsumerFunc, optionFuncs ...func(*storage.ConsumeOptions)) {
	if r.consumer == nil {
		// default connection ...
		consumer, err := r.newConsumer(ctx)
		if err != nil {
			glog.Error(ctx, "rabbitmq newConsumer error:", err)
			return
		}
		r.consumer = &consumer
	}
	options := storage.GetDefaultConsumeOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}
	// exchange exchangeType routingKey
	err := r.consumer.StartConsuming(ctx,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			glog.Debug(ctx, "Delivery:", d)
			m := new(Message)
			m.SetValues(map[string]interface{}{
				"body": string(d.Body),
			})
			m.SetRoutingKey(d.RoutingKey)
			m.SetId(d.MessageId)
			if d.Redelivered {
				m.SetErrorCount(d.DeliveryTag)
			}
			err := f(ctx, m)
			if err != nil {
				return rabbitmq.NackRequeue
			}
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		queueName,
		options.BindingRoutingKeys,
		rabbitmq.WithConsumeOptionsConsumerName(options.ConsumerName),
		rabbitmq.WithConsumeOptionsBindingExchangeName(options.BindingExchange.Name),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(options.BindingExchange.Kind),
		rabbitmq.WithConsumeOptionsConcurrency(options.Concurrency), // goroutine num
		rabbitmq.WithConsumeOptionsConsumerAutoAck(options.ConsumerAutoAck),
		rabbitmq.WithConsumeOptionsQOSPrefetch(options.QOSPrefetch),
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
