package rocketmq

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/task"
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
	mQueue := sdk.Runtime.GetRocketQueue(task.DefaultQueue) // get rabbitmq instance
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
					storage.WithRocketMqGroupName(spec.GroupName),
					storage.WithRocketMqAutoCommit(spec.AutoCommit),
					storage.WithRocketMqMaxReconsumeTimes(spec.MaxReconsumeTimes),
				)
			}
		}
		go mQueue.Run(ctx)
	} else {
		panic(gerror.New("sdk.Runtime.GetRocketQueue is nil!"))
	}
}
