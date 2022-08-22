package rocketmq

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
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
	Queue   storage.AdapterQueue
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
	t.Queue = nil // todo get rocketmq instance
	if t.Queue != nil {
		for _, worker := range t.Routers {
			spec := worker.GetSpec(ctx)
			if spec == nil {
				continue
			}
			// Consumer
			t.Queue.Consumer(ctx, spec.QueueName, worker.Handle)
		}
		go t.Queue.Run(ctx)
	} else {
		glog.Warning(ctx, "RocketMq is nil ...")
	}
}
