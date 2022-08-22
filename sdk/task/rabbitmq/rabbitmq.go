package rabbitmq

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/task"
)

const (
	SrvName     = "RabbitMqTask"
	QOSPrefetch = 100
)

var insRabbitMq = tRabbitMq{
	Routers: []task.RabbitMqTask{},
}

type tRabbitMq struct {
	vhost   string
	Queue   storage.AdapterQueue
	Routers []task.RabbitMqTask
}

func Service(vhost ...string) *tRabbitMq {
	if len(vhost) > 0 {
		insRabbitMq.vhost = vhost[0]
	}
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
	//t.Queue = sdk.Runtime.GetRabbitQueue("vhost") // get rabbitmq instance
	t.Queue = sdk.Runtime.GetRabbitQueue(t.vhost) // get rabbitmq instance
	if t.Queue != nil {
		for _, worker := range t.Routers {
			spec := worker.GetSpec(ctx)
			if spec == nil {
				continue
			}
			// Consumer
			go t.Queue.Consumer(ctx, spec.QueueName, worker.Handle,
				storage.WithConsumeOptionsBindingRoutingKeys(spec.GetRoutingKeys()),
				storage.WithConsumeOptionsBindingExchangeName(spec.Exchange),
				storage.WithConsumeOptionsBindingExchangeType(spec.ExchangeType),
				storage.WithConsumeOptionsConcurrency(spec.ConsumerNum),
				storage.WithConsumeOptionsConsumerName(spec.TaskName),
				storage.WithConsumeOptionsConsumerAutoAck(spec.AutoAck),
				storage.WithConsumeOptionsQOSPrefetch(QOSPrefetch),
			)
		}
	} else {
		glog.Warning(ctx, "RabbitMq is nil ...")
	}
}
