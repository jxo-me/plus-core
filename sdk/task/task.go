package task

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/storage"
)

// Handler 任务MQ路由的回调接口/**/
type Handler interface {
	Handle(ctx context.Context, msg storage.Messager) error
}

type ConsumerHandler interface {
	Handle(ctx context.Context, msg storage.Messager) (interface{}, error)
}

type Task interface {
	String() string
	Start(ctx context.Context)
}

type MemoryTask interface {
	GetSpec() *MemorySpec
	Handler
}

type MemorySpec struct {
	TaskName   string
	RoutingKey string
}

type RabbitMqTask interface {
	GetSpec(ctx context.Context) *RabbitMqSpec
	ConsumerHandler
}

type RabbitMqSpec struct {
	TaskName     string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	RoutingMap   map[string]ConsumerHandler
	ConsumerNum  int
	CTag         string
}

type RocketMqTask interface {
	GetSpec(ctx context.Context) *RocketMqSpec
	ConsumerHandler
}

type RocketMqSpec struct {
	TaskName     string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	RoutingMap   map[string]ConsumerHandler
	ConsumerNum  int
	CTag         string
}

type NsqTask interface {
	GetSpec(ctx context.Context) *NsqSpec
	ConsumerHandler
}

type NsqSpec struct {
	TaskName     string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	RoutingMap   map[string]ConsumerHandler
	ConsumerNum  int
	CTag         string
}

func WrapHandler(handler ConsumerHandler) storage.ConsumerFunc {
	return storage.ConsumerFunc(
		func(ctx context.Context, msg storage.Messager) error {
			glog.Debug(ctx, "handler result:", msg)
			_, err := handler.Handle(ctx, msg)
			if err != nil {
				glog.Error(ctx, "task handler error", err.Error())
			}
			return err
		},
	)
}

// CallbackFunc 消费结果统一回调
type CallbackFunc func(context.Context, interface{}) error

func CallbackWrapHandler(handler ConsumerHandler, callback CallbackFunc) storage.ConsumerFunc {
	return storage.ConsumerFunc(
		func(ctx context.Context, msg storage.Messager) error {
			glog.Debug(ctx, "handler result:", msg)
			data, err := handler.Handle(ctx, msg)
			if err != nil {
				glog.Error(ctx, "task handler error", err.Error())
			}
			err = callback(ctx, data)
			if err != nil {
				glog.Error(ctx, "task CallbackFunc error", err.Error())
				return err
			}
			return err
		},
	)
}
