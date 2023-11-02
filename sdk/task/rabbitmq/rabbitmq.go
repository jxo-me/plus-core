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
	gName := config.DefaultGroupName
	dQueue := sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RabbitmqQueueName, gName)) // get rabbitmq instance
	if dQueue == nil {
		panic(gerror.New("sdk.Runtime.GetRabbitQueue default group is nil!"))
	}
	t.Queue[gName] = dQueue
	q = dQueue
	// register task
	for _, worker := range t.Routers {
		spec := worker.GetSpec(ctx)
		if spec == nil {
			glog.Warning(ctx, "get tRabbitMq spec is nil ignore...")
			continue
		}
		if spec.Vhost != "" {
			gName = spec.Vhost
			if cQueue, ok := t.Queue[gName]; ok {
				t.Queue[gName] = cQueue
				q = cQueue
			} else {
				// get custom queue
				cQueue = sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RabbitmqQueueName, gName)) // get rabbitmq instance
				if cQueue != nil {
					t.Queue[gName] = cQueue
					q = cQueue
				} else {
					glog.Warning(ctx, fmt.Sprintf("task name: %s, get queue %s group is nil, use default queue.", spec.TaskName, gName))
				}
			}
			// use default queue
		}

		glog.Info(ctx, fmt.Sprintf("task name: %s, use queue group name: %s", spec.TaskName, gName))
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
				queue.WithRabbitMqConsumeOptionsExchangePassive(spec.Passive),
				queue.WithRabbitMqConsumeOptionsExchangeDeclare(spec.Declare),
				queue.WithRabbitMqConsumeOptionsExchangeDurable(spec.Durable),
			)
		}
	}
}
