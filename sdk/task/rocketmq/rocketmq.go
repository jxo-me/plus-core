package rocketmq

import (
	"context"
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
}

type tRocketMq struct {
	Routers []task.RocketMqTask
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
	glog.Info(ctx, "RocketMq task start ...")
	mQueue := sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RocketQueueName, config.DefaultGroupName)) // get rabbitmq instance
	if mQueue != nil {
		for _, worker := range t.Routers {
			spec := worker.GetSpec(ctx)
			if spec == nil {
				glog.Warning(ctx, "get tRocketMq spec is nil ignore...")
				continue
			}
			for i := 0; i < spec.ConsumerNum; i++ {
				// Consumer
				mQueue.Consumer(ctx, spec.TopicName, worker.Handle,
					queue.WithRocketMqGroupName(spec.GroupName),
					queue.WithRocketMqAutoCommit(spec.AutoCommit),
					queue.WithRocketMqMaxReconsumeTimes(spec.MaxReconsumeTimes),
				)
			}
		}
		go mQueue.Run(ctx)
	} else {
		panic(gerror.New("sdk.Runtime.GetRocketQueue is nil!"))
	}
}
