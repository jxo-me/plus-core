package task

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/storage"
)

type MemoryTask interface {
	GetSpec() *MemorySpec
	Handle(ctx context.Context, msg storage.Messager) error
}

type MemorySpec struct {
	TaskName   string
	RoutingKey string
}

type RabbitMqTask interface {
	GetSpec(ctx context.Context) *RabbitMqSpec
}

type RabbitMqSpec struct {
	TaskName     string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	RoutingMap   map[string]Consumer
	ConsumerNum  int
	CTag         string
}

type RocketMqTask interface {
	GetSpec(ctx context.Context) *RocketMqSpec
}

type RocketMqSpec struct {
	TaskName     string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	RoutingMap   map[string]Consumer
	ConsumerNum  int
	CTag         string
}

// ConsumePart /*
type ConsumePart interface {
	Route(routingKey string) (handler Consumer, ifExist bool)
}

// Consumer 任务MQ路由的回调接口/**/
type Consumer interface {
	Handle(body []byte) error
}

// ConsumerFunc 将一个符合签名要求的函数转换成 Consumer 接口/*
type ConsumerFunc func([]byte) error

func (f ConsumerFunc) Handle(body []byte) error {
	return f(body)
}

type ConsumerHandler interface {
	Handle(ctx context.Context, body []byte) (interface{}, error)
}

func WrapHandler(ctx context.Context, handler ConsumerHandler) Consumer {
	return ConsumerFunc(
		func(body []byte) error {
			glog.Debug(ctx, "handler result:", string(body))
			_, err := handler.Handle(ctx, body)
			if err != nil {
				glog.Error(ctx, "rabbitmq task handler error", err.Error())
			}
			return nil
		},
	)

}
