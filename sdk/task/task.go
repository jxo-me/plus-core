package task

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/storage"
)

const (
	DefaultQueue = "default"
)

// Handler 任务MQ路由的回调接口/**/
type Handler interface {
	Handle(ctx context.Context, msg storage.Messager) error
}

type SubTask interface {
	RoutingKey() string
	Handle(ctx context.Context, msg storage.Messager) (interface{}, error)
}

type Service interface {
	String() string
	Start(ctx context.Context)
}

type TasksService interface {
	Service
	AddServices(services ...Service) TasksService
}

type RabbitMqService interface {
	Service
	AddTasks(task ...RabbitMqTask) RabbitMqService
}

type RocketMqService interface {
	Service
	AddTasks(task ...RocketMqTask) RocketMqService
}

type MemoryService interface {
	Service
	AddTasks(task ...MemoryTask) MemoryService
}

type NsqService interface {
	Service
	AddTasks(task ...NsqTask) NsqService
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
	Handler
}

type RabbitMqSpec struct {
	TaskName     string
	Vhost        string
	RoutingKeys  []string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	SubTasks     []SubTask
	ConsumerNum  int
	CoroutineNum int
	Prefetch     int
	AutoAck      bool
}

func (r *RabbitMqSpec) GetRoutingKeys() []string {
	r.RoutingKeys = make([]string, 0)
	r.RoutingKeys = append(r.RoutingKeys, r.RoutingKey)
	for _, subHandler := range r.SubTasks {
		r.RoutingKeys = append(r.RoutingKeys, subHandler.RoutingKey())
	}

	return r.RoutingKeys
}

func (r *RabbitMqSpec) Route(routingKey string) (handler SubTask, ifExist bool) {
	for _, subHandler := range r.SubTasks {
		if subHandler.RoutingKey() == routingKey {
			return subHandler, true
		}
	}
	return nil, false
}

type RocketMqTask interface {
	GetSpec(ctx context.Context) *RocketMqSpec
	Handler
}

type RocketMqSpec struct {
	TaskName          string
	GroupName         string
	TopicName         string
	SubTasks          []SubTask
	ConsumerNum       int
	MaxReconsumeTimes int32
	AutoCommit        bool
}

func (r *RocketMqSpec) Route(routingKey string) (handler SubTask, ifExist bool) {
	for _, subHandler := range r.SubTasks {
		if subHandler.RoutingKey() == routingKey {
			return subHandler, true
		}
	}
	return nil, false
}

type NsqTask interface {
	GetSpec(ctx context.Context) *NsqSpec
	Handler
}

type NsqSpec struct {
	TaskName     string
	RoutingKey   string
	Exchange     string
	ExchangeType string
	QueueName    string
	RoutingMap   map[string]SubTask
	ConsumerNum  int
	CTag         string
}

func WrapHandler(handler SubTask) storage.ConsumerFunc {
	return storage.ConsumerFunc(
		func(ctx context.Context, msg storage.Messager) error {
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

func CallbackWrapHandler(handler SubTask, callback CallbackFunc) storage.ConsumerFunc {
	return storage.ConsumerFunc(
		func(ctx context.Context, msg storage.Messager) error {
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
