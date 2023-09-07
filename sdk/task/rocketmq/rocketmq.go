package rocketmq

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/core/v2/task"
	"github.com/jxo-me/plus-core/sdk/v2"
	"github.com/jxo-me/plus-core/sdk/v2/config"
)

const (
	SrvName = "RocketMqTask"
)

var insRocketmq = tRocketMq{
	Routers: []task.RocketMqTask{},
	Queue:   map[string]queue.IQueue{},
}

type tRocketMq struct {
	Routers []task.RocketMqTask
	Queue   map[string]queue.IQueue
}

func Service() *tRocketMq {
	return &insRocketmq
}

func (t *tRocketMq) String() string {
	return SrvName
}

func (t *tRocketMq) AddTasks(tasks ...task.RocketMqTask) task.RocketMqService {
	t.Routers = tasks
	return t
}

func (t *tRocketMq) Start(ctx context.Context) {
	var q queue.IQueue
	glog.Info(ctx, "RocketMq task start ...")
	gName := config.DefaultGroupName
	dQueue := sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RocketQueueName, gName)) // get rocketmq instance
	if dQueue == nil {
		panic(gerror.New("sdk.Runtime.GetRocketQueue default group is nil!"))
	}
	t.Queue[gName] = dQueue
	q = dQueue
	for _, worker := range t.Routers {
		spec := worker.GetSpec(ctx)
		if spec == nil {
			glog.Warning(ctx, "get tRocketMq spec is nil ignore...")
			continue
		}
		if spec.GroupName != "" {
			gName = spec.GroupName
			if cQueue, ok := t.Queue[gName]; ok {
				t.Queue[gName] = cQueue
				q = cQueue
			} else {
				// get custom queue
				cQueue = sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RocketQueueName, gName)) // get rabbitmq instance
				if cQueue != nil {
					t.Queue[gName] = cQueue
					q = cQueue
				} else {
					glog.Warning(ctx, fmt.Sprintf("task name: %s, get queue %s group is nil, use default queue.", spec.TaskName, gName))
				}
			}
			// use default queue
		}
		glog.Info(ctx, fmt.Sprintf("rocketmq task name: %s, use queue group name: %s", spec.TaskName, gName))
		for i := 0; i < spec.ConsumerNum; i++ {
			// Consumer
			q.Consumer(ctx, spec.TopicName, worker.Handle,
				queue.WithRocketMqGroupName(spec.GroupName),
				queue.WithRocketMqAutoCommit(spec.AutoCommit),
				queue.WithRocketMqMaxReconsumeTimes(spec.MaxReconsumeTimes),
			)
		}
	}
	go q.Run(ctx)
}
