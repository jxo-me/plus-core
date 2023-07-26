package task

import (
	"context"
)

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

type NsqTask interface {
	GetSpec(ctx context.Context) *NsqSpec
	IHandler
}

type NsqService interface {
	IService
	AddTasks(task ...NsqTask) NsqService
}
