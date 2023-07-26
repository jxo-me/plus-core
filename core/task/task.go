package task

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/message"
	"github.com/jxo-me/plus-core/core/queue"
)

const (
	DefaultQueue = "default"
)

// IHandler 任务MQ路由的回调接口/**/
type IHandler interface {
	Handle(ctx context.Context, msg message.IMessage) error
}

type SubTask interface {
	RoutingKey() string
	Handle(ctx context.Context, msg message.IMessage) (interface{}, error)
}

type TasksService interface {
	IService
	AddServices(services ...IService) TasksService
}

func WrapHandler(handler SubTask) queue.ConsumerFunc {
	return queue.ConsumerFunc(
		func(ctx context.Context, msg message.IMessage) error {
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

func CallbackWrapHandler(handler SubTask, callback CallbackFunc) queue.ConsumerFunc {
	return queue.ConsumerFunc(
		func(ctx context.Context, msg message.IMessage) error {
			data, err := handler.Handle(ctx, msg)
			if err != nil {
				glog.Error(ctx, "task handler error", err.Error())
			}
			err = callback(ctx, data)
			if err != nil {
				glog.Error(ctx, "task CallbackFunc error", err.Error())
			}
			return err
		},
	)
}
