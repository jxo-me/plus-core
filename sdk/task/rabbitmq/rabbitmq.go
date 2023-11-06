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
		if spec.Group != "" {
			gName = spec.Group
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
					glog.Warning(ctx, fmt.Sprintf("task name: %s, get queue %s config group is nil, use default config group. Exchange: %s ExchangeType: %s QueueName: %s", spec.TaskName, gName, spec.Exchange, spec.ExchangeType, spec.QueueName))
				}
			}
			// use default queue
		}

		glog.Info(ctx, fmt.Sprintf("task name: %s, use config group name: %s, Exchange: %s ExchangeType: %s QueueName: %s", spec.TaskName, gName, spec.Exchange, spec.ExchangeType, spec.QueueName))
		for i := 0; i < spec.ConsumerNum; i++ {
			var optFuncs []func(*queue.ConsumeOptions)
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsBindingRoutingKeys(spec.GetRoutingKeys()))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsBindingExchangeName(spec.Exchange))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsBindingExchangeType(spec.ExchangeType))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsConcurrency(spec.CoroutineNum))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsConsumerName(fmt.Sprintf("%s.%02d", spec.TaskName, i+1)))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsConsumerAutoAck(spec.AutoAck))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsQOSPrefetch(spec.Prefetch))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsExchangePassive(spec.Passive))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsExchangeDeclare(spec.Declare))
			optFuncs = append(optFuncs, queue.WithRabbitMqConsumeOptionsExchangeDurable(spec.Durable))
			// Consumer
			q.Consumer(ctx, spec.QueueName, worker.Handle, optFuncs...)
		}
	}
}
