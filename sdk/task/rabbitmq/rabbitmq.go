package rabbitmq

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/task"
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
	for _, worker := range t.Routers {
		spec := worker.GetSpec(ctx)
		if spec == nil {
			continue
		}
		mQueue := sdk.Runtime.GetRabbitQueue(task.DefaultQueue) // get rabbitmq instance
		if mQueue != nil {
			for i := 0; i < spec.ConsumerNum; i++ {
				// Consumer
				go mQueue.Consumer(ctx, spec.QueueName, worker.Handle,
					storage.WithRabbitMqConsumeOptionsBindingRoutingKeys(spec.GetRoutingKeys()),
					storage.WithRabbitMqConsumeOptionsBindingExchangeName(spec.Exchange),
					storage.WithRabbitMqConsumeOptionsBindingExchangeType(spec.ExchangeType),
					storage.WithRabbitMqConsumeOptionsConcurrency(spec.CoroutineNum),
					storage.WithRabbitMqConsumeOptionsConsumerName(fmt.Sprintf("%s.%02d", spec.TaskName, i+1)),
					storage.WithRabbitMqConsumeOptionsConsumerAutoAck(spec.AutoAck),
					storage.WithRabbitMqConsumeOptionsQOSPrefetch(spec.Prefetch),
				)
			}
		} else {
			panic(gerror.New("sdk.Runtime.GetRabbitQueue is nil!"))
		}
	}
}
