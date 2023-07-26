package task

import (
	"context"
)

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

type RocketMqTask interface {
	GetSpec(ctx context.Context) *RocketMqSpec
	IHandler
}

type RocketMqService interface {
	IService
	AddTasks(task ...RocketMqTask) RocketMqService
}
