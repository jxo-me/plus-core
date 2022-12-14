package nsq

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/task"
)

const (
	SrvName = "NsqTask"
)

var insNsq = tNsq{
	Routers: []task.NsqTask{},
}

type tNsq struct {
	Queue   storage.AdapterQueue
	Routers []task.NsqTask
}

func Service() *tNsq {
	return &insNsq
}

func (t *tNsq) String() string {
	return SrvName
}

func (t *tNsq) AddTasks(tasks ...task.NsqTask) task.NsqService {
	t.Routers = tasks
	return t
}

func (t *tNsq) Start(ctx context.Context) {
	glog.Info(ctx, "RocketMq task start ...")
	t.Queue = nil // todo get rocketmq instance
	if t.Queue != nil {
		for _, worker := range t.Routers {
			spec := worker.GetSpec(ctx)
			if spec == nil {
				continue
			}
			// Consumer
			go t.Queue.Consumer(ctx, spec.QueueName, worker.Handle)
		}
	} else {
		glog.Warning(ctx, "RocketMq is nil ...")
	}
}
