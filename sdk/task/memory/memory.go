package memory

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/task"
)

const (
	SrvName      = "MemoryTask"
	DefaultQueue = "default"
)

var instMemory = tMemory{
	Routers: []task.MemoryTask{},
}

type tMemory struct {
	Queue   storage.AdapterQueue
	Routers []task.MemoryTask
}

func Service() *tMemory {
	return &instMemory
}

func (t *tMemory) String() string {
	return SrvName
}

func (t *tMemory) AddTasks(tasks ...task.MemoryTask) task.MemoryService {
	t.Routers = tasks
	return t
}

func (t *tMemory) Start(ctx context.Context) {
	glog.Info(ctx, "MemoryMq task start ...")
	t.Queue = sdk.Runtime.GetMemoryQueue(DefaultQueue)
	if t.Queue != nil {
		for _, worker := range t.Routers {
			sp := worker.GetSpec()
			t.Queue.Consumer(ctx, sp.RoutingKey, worker.Handle)
		}
		go t.Queue.Run(ctx)
	} else {
		glog.Warning(ctx, "MemoryMq is nil ...")
	}
}
