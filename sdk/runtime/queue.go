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

// Register 注册消费者
func (e *Queue) Register(ctx context.Context, name string, f storage.ConsumerFunc) {
	e.queue.Register(ctx, name, f)
}

// Append 增加数据到生产者
func (e *Queue) Append(ctx context.Context, message storage.Messager) error {
	values := message.GetValues()
	if values == nil {
		values = make(map[string]interface{})
	}
	values[storage.PrefixKey] = e.prefix
	return e.queue.Append(ctx, message)
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
