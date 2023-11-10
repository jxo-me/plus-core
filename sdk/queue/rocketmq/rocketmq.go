package rocketmq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/gogf/gf/v2/os/glog"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/message"
	"sync"
)

func NewRocketMQ(
	ctx context.Context,
	urls []string,
	credentials *primitive.Credentials,
	logger *glog.Logger,
) (*RocketMQ, error) {
	r := &RocketMQ{
		Urls:        urls,
		Credentials: credentials,
		consumers:   map[string]rocketmq.PushConsumer{},
		producers:   map[string]rocketmq.Producer{},
	}

	rlog.SetLogger(&QLoger{Logger: logger})
	return r, nil
}

// RocketMQ cache implement
type RocketMQ struct {
	Urls        []string
	consumers   map[string]rocketmq.PushConsumer
	producers   map[string]rocketmq.Producer
	mux         sync.RWMutex
	Credentials *primitive.Credentials
}

func (r *RocketMQ) String() string {
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

func (r *RocketMQ) newConsumer(ctx context.Context, opt queueLib.ConsumeOptions) (rocketmq.PushConsumer, error) {
	if r.Credentials != nil {
		return rocketmq.NewPushConsumer(
			consumer.WithGroupName(opt.GroupName),
			consumer.WithNsResolver(primitive.NewPassthroughResolver(r.Urls)),
			consumer.WithConsumerModel(consumer.Clustering),
			consumer.WithMaxReconsumeTimes(opt.MaxReconsumeTimes),
			consumer.WithAutoCommit(opt.AutoCommit),
			consumer.WithCredentials(*r.Credentials),
			//consumer.WithCredentials(primitive.Credentials{
			//	AccessKey: "RocketMQ",
			//	SecretKey: "12345678",
			//}),
		)
	}
	return rocketmq.NewPushConsumer(
		consumer.WithGroupName(opt.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(r.Urls)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithMaxReconsumeTimes(opt.MaxReconsumeTimes),
		consumer.WithAutoCommit(opt.AutoCommit),
		//consumer.WithCredentials(*r.Credentials),
		//consumer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
}

func (r *RocketMQ) newProducer(ctx context.Context, opt queueLib.PublishOptions) (rocketmq.Producer, error) {
	if r.Credentials != nil {
		return rocketmq.NewProducer(
			producer.WithNsResolver(primitive.NewPassthroughResolver(r.Urls)),
			producer.WithRetry(opt.RetryTimes),
			producer.WithCredentials(*r.Credentials),
			//producer.WithCredentials(primitive.Credentials{
			//	AccessKey: "RocketMQ",
			//	SecretKey: "12345678",
			//}),
			//producer.WithNamespace("namespace"),
		)
	}

	return rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(r.Urls)),
		producer.WithRetry(opt.RetryTimes),
		//producer.WithCredentials(*r.Credentials),
		//producer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
		//producer.WithNamespace("namespace"),
	)
}

// Publish 消息入生产者
func (r *RocketMQ) Publish(ctx context.Context, message messageLib.IMessage, optionFuncs ...func(*queueLib.PublishOptions)) error {
	options := queueLib.PublishOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}
	var p rocketmq.Producer
	var err error
	var ok bool
	if p, ok = r.producers[options.GroupName]; !ok {
		p, err = r.newProducer(ctx, options)
		if err != nil {
			glog.Error(ctx, "RocketMQ newConsumer error:", err)
			return err
		}
		r.producers[options.GroupName] = p
	}
	// encode message
	rb, err := json.Marshal(message.GetValue())
	if err != nil {
		return err
	}
	_, err = p.SendSync(
		ctx,
		&primitive.Message{
			Topic: message.GetRoutingKey(),
			Body:  rb,
		})
	return err
}

func (r *RocketMQ) RpcRequest(ctx context.Context, key string, data []byte, optionFuncs ...func(*queueLib.ClientOptions)) ([]byte, error) {
	return nil, nil
}

// Consumer 监听消费者
func (r *RocketMQ) Consumer(ctx context.Context, topicName string, f queueLib.ConsumerFunc, optionFuncs ...func(*queueLib.ConsumeOptions)) {
	options := queueLib.GetDefaultConsumeOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}
	var c rocketmq.PushConsumer
	var err error
	var ok bool
	r.mux.Lock()
	defer r.mux.Unlock()
	if c, ok = r.consumers[options.GroupName]; !ok {
		c, err = r.newConsumer(ctx, options)
		if err != nil {
			glog.Error(ctx, "RocketMQ newConsumer error:", err)
			return
		}
		r.consumers[options.GroupName] = c
	}
	err = c.Subscribe(
		topicName,
		consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {
				if len(msgs[i].Body) > 0 {
					glog.Debugf(ctx, "rocketmq consumed: %v\n", msgs[i])
					m := new(message.Message)
					m.SetValue(msgs[i].Body)
					m.SetRoutingKey(msgs[i].GetTags())
					m.SetId(msgs[i].MsgId)
					m.SetErrorCount(uint64(msgs[i].ReconsumeTimes))
					err = f(ctx, nil, m)
					if err != nil {
						glog.Warning(ctx, "RocketMQ Rollback msg:", m)
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
	for _, pushConsumer := range r.consumers {
		err := pushConsumer.Start()
		if err != nil {
			glog.Warning(ctx, "rocketmq consumer Start error", err)
			continue
		}
	}
	return
}

func (r *RocketMQ) Shutdown(ctx context.Context) {
	for _, pd := range r.producers {
		err := pd.Shutdown()
		if err != nil {
			glog.Warning(ctx, "rocketmq producer Close error", err)
		}
	}
	for _, pushConsumer := range r.consumers {
		err := pushConsumer.Shutdown()
		if err != nil {
			glog.Warning(ctx, "rocketmq consumer Close error", err)
		}
	}
}
