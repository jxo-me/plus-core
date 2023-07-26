package task

import (
	"context"
)

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
	if r.RoutingKey != "" {
		r.RoutingKeys = append(r.RoutingKeys, r.RoutingKey)
	}
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

type RabbitMqTask interface {
	GetSpec(ctx context.Context) *RabbitMqSpec
	IHandler
}

type RabbitMqService interface {
	IService
	AddTasks(task ...RabbitMqTask) RabbitMqService
}
