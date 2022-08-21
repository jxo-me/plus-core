package queue

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jxo-me/plus-core/sdk/storage"
)

func NewRocketMQ(
	ctx context.Context,
	urls []string,
	consumerOptions *RocketConsumerOptions,
	producerOptions *RocketProducerOptions,
	credentials *primitive.Credentials,
) (*RocketMQ, error) {
	var err error
	r := &RocketMQ{
		Urls:            urls,
		ConsumerOptions: consumerOptions,
		ProducerOptions: producerOptions,
		Credentials:     *credentials,
	}
	if credentials != nil {
		r.Credentials = *credentials
	}
	r.consumer, err = r.newConsumer(ctx, r.ConsumerOptions)
	if err != nil {
		return nil, err
	}
	r.producer, err = r.newProducer(ctx, r.ProducerOptions)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// RocketMQ cache implement
type RocketMQ struct {
	Urls            []string
	consumer        rocketmq.PushConsumer
	ConsumerOptions *RocketConsumerOptions
	producer        rocketmq.Producer
	ProducerOptions *RocketProducerOptions
	Credentials     primitive.Credentials
}

func (RocketMQ) String() string {
	return "rocketmq"
}

type RocketConsumerOptions struct {
	GroupName         string
	MaxReconsumeTimes int32
}

type RocketProducerOptions struct {
	GroupName  string
	RetryTimes int
}

func (r *RocketMQ) newConsumer(ctx context.Context, options *RocketConsumerOptions) (rocketmq.PushConsumer, error) {
	if options == nil {
		r.ConsumerOptions = &RocketConsumerOptions{
			GroupName:         "DEFAULT_CONSUMER",
			MaxReconsumeTimes: -1,
		}
	}
	return rocketmq.NewPushConsumer(
		consumer.WithGroupName(r.ConsumerOptions.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(r.Urls)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithMaxReconsumeTimes(r.ConsumerOptions.MaxReconsumeTimes),
		consumer.WithAutoCommit(false),
		consumer.WithCredentials(r.Credentials),
		//consumer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
}

func (r *RocketMQ) newProducer(ctx context.Context, options *RocketProducerOptions) (rocketmq.Producer, error) {
	if options == nil {
		r.ProducerOptions = &RocketProducerOptions{
			GroupName:  "DEFAULT_CONSUMER",
			RetryTimes: 3,
		}
	}
	return rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(r.Urls)),
		producer.WithRetry(r.ProducerOptions.RetryTimes),
		producer.WithCredentials(r.Credentials),
		//producer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
		//producer.WithNamespace("namespace"),
	)
}

// Publish 消息入生产者
func (r *RocketMQ) Publish(ctx context.Context, message storage.Messager, optionFuncs ...func(*storage.PublishOptions)) error {
	//
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	_, err = r.producer.SendSync(
		ctx,
		&primitive.Message{
			Topic: message.GetRoutingKey(),
			Body:  rb,
		})
	return err
}

// Consumer 监听消费者
func (r *RocketMQ) Consumer(ctx context.Context, topicName string, f storage.ConsumerFunc, optionFuncs ...func(*storage.ConsumeOptions)) {
	err := r.consumer.Subscribe(
		topicName,
		consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {
				if len(msgs[i].Body) > 0 {
					glog.Debugf(ctx, "rocketmq consumed: %s\n", string(msgs[i].Body))
					m := new(Message)
					m.SetValues(gconv.Map(msgs[i].Body))
					m.SetRoutingKey(msgs[i].GetTags())
					m.SetId(msgs[i].MsgId)
					err := f(ctx, m)
					if err != nil {
						return consumer.Rollback, err
					}
				}
			}
			return consumer.ConsumeSuccess, nil
		},
	)
	if err != nil {
		glog.Errorf(ctx, "rocketmq consumer Subscribe error:%v", err)
		return
	}
}

func (r *RocketMQ) Run(ctx context.Context) {
	err := r.consumer.Start()
	if err != nil {
		glog.Warning(ctx, "rocketmq consumer Start error", err)
	}
	return
}

func (r *RocketMQ) Shutdown(ctx context.Context) {
	if r.producer != nil {
		err := r.producer.Shutdown()
		if err != nil {
			glog.Warning(ctx, "rocketmq producer Close error", err)
		}
	}
	if r.consumer != nil {
		err := r.consumer.Shutdown()
		if err != nil {
			glog.Warning(ctx, "rocketmq consumer Close error", err)
		}
	}
}
