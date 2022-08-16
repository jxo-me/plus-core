package runtime

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
)

// NewQueue 创建对应上下文队列
func NewQueue(prefix string, queue storage.AdapterQueue) storage.AdapterQueue {
	return &Queue{
		prefix: prefix,
		queue:  queue,
	}
}

type Queue struct {
	prefix string
	queue  storage.AdapterQueue
}

func (e *Queue) String() string {
	return e.queue.String()
}

// Consumer 注册消费者
func (e *Queue) Consumer(ctx context.Context, name string, f storage.ConsumerFunc) {
	e.queue.Consumer(ctx, name, f)
}

// Publish 数据生产者
func (e *Queue) Publish(ctx context.Context, message storage.Messager) error {
	values := message.GetValues()
	if values == nil {
		values = make(map[string]interface{})
	}
	values[storage.PrefixKey] = e.prefix
	return e.queue.Publish(ctx, message)
}

// Run 运行
func (e *Queue) Run(ctx context.Context) {
	e.queue.Run(ctx)
}

// Shutdown 停止
func (e *Queue) Shutdown(ctx context.Context) {
	if e.queue != nil {
		e.queue.Shutdown(ctx)
	}
}
