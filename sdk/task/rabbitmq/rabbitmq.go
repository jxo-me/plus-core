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
}

type tRabbitMq struct {
	Routers []task.RabbitMqTask
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
	glog.Info(ctx, "RabbitMq task start ...")
	mQueue := sdk.Runtime.QueueRegistry().Get(config.GetQueueName(config.RabbitmqQueueName, config.DefaultGroupName)) // get rabbitmq instance
	if mQueue != nil {
		for _, worker := range t.Routers {
			spec := worker.GetSpec(ctx)
			if spec == nil {
				glog.Warning(ctx, "get tRabbitMq spec is nil ignore...")
				continue
			}
			for i := 0; i < spec.ConsumerNum; i++ {
				// Consumer
				mQueue.Consumer(ctx, spec.QueueName, worker.Handle,
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
	} else {
		panic(gerror.New("sdk.Runtime.GetRabbitQueue is nil!"))
	}
}
