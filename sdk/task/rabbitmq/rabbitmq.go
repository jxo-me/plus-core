package rabbitmq

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
	SrvName = "RabbitMqTask"
)

var insRabbitMq = tRabbitMq{
	Routers: []task.RabbitMqTask{},
	Queue:   map[string]queue.IQueue{},
}

type tRabbitMq struct {
	Routers []task.RabbitMqTask
	Queue   map[string]queue.IQueue
}

func Service() *tRabbitMq {
	return &insRabbitMq
}

func (t *tRabbitMq) String() string {
	return SrvName
}

func (t *tRabbitMq) AddTasks(tasks ...task.RabbitMqTask) task.RabbitMqService {
	t.Routers = tasks
	return t
}

func (t *tRabbitMq) Start(ctx context.Context) {
	var q queue.IQueue
	glog.Info(ctx, "RabbitMq task start ...")
	dQueue := sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RabbitmqQueueName, config.DefaultGroupName)) // get rabbitmq instance
	if dQueue != nil {
		t.Queue[config.DefaultGroupName] = dQueue
	} else {
		panic(gerror.New("sdk.Runtime.GetRabbitQueue is nil!"))
	}
	// register task
	for _, worker := range t.Routers {
		spec := worker.GetSpec(ctx)
		if spec == nil {
			glog.Warning(ctx, "get tRabbitMq spec is nil ignore...")
			continue
		}
		q = dQueue
		gName := spec.Vhost
		if gName != "" {
			if cQueue, ok := t.Queue[gName]; ok {
				if cQueue != nil {
					t.Queue[gName] = cQueue
					q = cQueue
				} else {
					glog.Warning(ctx, fmt.Sprintf("cQueue %s is nil use default queue: %s", gName, q.String()))
				}
			} else {
				// get custom queue
				cQueue = sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RabbitmqQueueName, gName)) // get rabbitmq instance
				if cQueue != nil {
					t.Queue[gName] = cQueue
					q = cQueue
				} else {
					glog.Warning(ctx, fmt.Sprintf("get cQueue %s is nil use default queue: %s", gName, q.String()))
				}
			}
			glog.Info(ctx, fmt.Sprintf("queue group name: %s queue name: %s", gName, q.String()))
			// use default queue
		}

		for i := 0; i < spec.ConsumerNum; i++ {
			// Consumer
			q.Consumer(ctx, spec.QueueName, worker.Handle,
				queue.WithRabbitMqConsumeOptionsBindingRoutingKeys(spec.GetRoutingKeys()),
				queue.WithRabbitMqConsumeOptionsBindingExchangeName(spec.Exchange),
				queue.WithRabbitMqConsumeOptionsBindingExchangeType(spec.ExchangeType),
				queue.WithRabbitMqConsumeOptionsConcurrency(spec.CoroutineNum),
				queue.WithRabbitMqConsumeOptionsConsumerName(fmt.Sprintf("%s.%02d", spec.TaskName, i+1)),
				queue.WithRabbitMqConsumeOptionsConsumerAutoAck(spec.AutoAck),
				queue.WithRabbitMqConsumeOptionsQOSPrefetch(spec.Prefetch),
			)
		}
	}
}
