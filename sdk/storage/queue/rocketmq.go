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

// RocketMQ cache implement
type RocketMQ struct {
	Url             []string
	consumer        rocketmq.PushConsumer
	consumerOptions *ConsumerOptions
	producer        rocketmq.Producer
	producerOptions *ProducerOptions
	Credentials     primitive.Credentials
}

func (RocketMQ) String() string {
	return "rocketmq"
}

type ConsumerOptions struct {
	GroupName         string
	MaxReconsumeTimes int32
}

type ProducerOptions struct {
	GroupName  string
	RetryTimes int
}

func (r *RocketMQ) newConsumer(ctx context.Context, options *ConsumerOptions) (rocketmq.PushConsumer, error) {
	if options == nil {
		r.consumerOptions = &ConsumerOptions{
			GroupName:         "DEFAULT_CONSUMER",
			MaxReconsumeTimes: -1,
		}
	}
	return rocketmq.NewPushConsumer(
		//consumer.WithGroupName("GAME_RECORD"),
		consumer.WithGroupName(r.consumerOptions.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(r.Url)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithMaxReconsumeTimes(r.consumerOptions.MaxReconsumeTimes),
		consumer.WithCredentials(r.Credentials),
		//consumer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
}

func (r *RocketMQ) newProducer(ctx context.Context, options *ProducerOptions) (rocketmq.Producer, error) {
	if options == nil {
		r.producerOptions = &ProducerOptions{
			GroupName:  "DEFAULT_CONSUMER",
			RetryTimes: 3,
		}
	}
	return rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(r.Url)),
		producer.WithRetry(r.producerOptions.RetryTimes),
		producer.WithCredentials(r.Credentials),
		//producer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
		//producer.WithNamespace("namespace"),
	)
}

// Publish 消息入生产者
func (r *RocketMQ) Publish(ctx context.Context, message storage.Messager) error {
	//
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	_, err = r.producer.SendSync(
		ctx,
		&primitive.Message{
			Topic: message.GetStream(),
			Body:  rb,
		})
	return err
}

// Consumer 监听消费者
func (r *RocketMQ) Consumer(ctx context.Context, name string, f storage.ConsumerFunc) {
	err := r.consumer.Subscribe(
		name,
		consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {
				if len(msgs[i].Body) > 0 {
					glog.Printf(ctx, "rocketmq consumed: %s\n", string(msgs[i].Body))
					m := new(Message)
					m.SetValues(gconv.Map(msgs[i].Body))
					m.SetStream(msgs[i].GetTags())
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
