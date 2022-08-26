package rabbitmq

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk"
	"github.com/jxo-me/plus-core/sdk/config"
	"github.com/jxo-me/plus-core/sdk/storage"
	"github.com/jxo-me/plus-core/sdk/storage/queue"
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
		for i := 0; i < spec.ConsumerNum; i++ {
			mQueue := sdk.Runtime.GetRabbitQueue(spec.Vhost) // get rabbitmq instance
			if mQueue != nil {
				// Consumer
				go mQueue.Consumer(ctx, spec.QueueName, worker.Handle,
					storage.WithConsumeOptionsBindingRoutingKeys(spec.GetRoutingKeys()),
					storage.WithConsumeOptionsBindingExchangeName(spec.Exchange),
					storage.WithConsumeOptionsBindingExchangeType(spec.ExchangeType),
					storage.WithConsumeOptionsConcurrency(spec.CoroutineNum),
					storage.WithConsumeOptionsConsumerName(fmt.Sprintf("%s.%02d", spec.TaskName, i+1)),
					storage.WithConsumeOptionsConsumerAutoAck(spec.AutoAck),
					storage.WithConsumeOptionsQOSPrefetch(QOSPrefetch),
				)
			} else {
				glog.Warning(ctx, "RabbitMq is nil ... NewRabbitMQ ...")
				cfg := config.QueueRabbit().GetCfg()
				// Use the vhost defined in the task first
				cfg.Vhost = spec.Vhost
				// get config connection
				mq, err := queue.NewRabbitMQ(ctx, config.QueueRabbit().GetDsn(), config.QueueRabbit().GetReconnectInterval(), cfg)
				if err != nil {
					glog.Error(ctx, "task NewRabbitMQ error:", err)
				}
				sdk.Runtime.SetQueueAdapter(spec.Vhost, mq)
				// Consumer
				go mq.Consumer(ctx, spec.QueueName, worker.Handle,
					storage.WithConsumeOptionsBindingRoutingKeys(spec.GetRoutingKeys()),
					storage.WithConsumeOptionsBindingExchangeName(spec.Exchange),
					storage.WithConsumeOptionsBindingExchangeType(spec.ExchangeType),
					storage.WithConsumeOptionsConcurrency(spec.CoroutineNum),
					storage.WithConsumeOptionsConsumerName(fmt.Sprintf("%s.%02d", spec.TaskName, i+1)),
					storage.WithConsumeOptionsConsumerAutoAck(spec.AutoAck),
					storage.WithConsumeOptionsQOSPrefetch(QOSPrefetch),
				)
			}
		}
	}
}
