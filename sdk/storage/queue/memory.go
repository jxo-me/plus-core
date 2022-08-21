package queue

import (
	"context"
	"github.com/google/uuid"
	"github.com/jxo-me/plus-core/sdk/storage"
	"sync"
)

type queue chan storage.Messager

// NewMemory 内存模式
func NewMemory(poolNum uint) *Memory {
	return &Memory{
		queue:   new(sync.Map),
		PoolNum: poolNum,
	}
}

type Memory struct {
	queue   *sync.Map
	wait    sync.WaitGroup
	mutex   sync.RWMutex
	PoolNum uint
}

func (*Memory) String() string {
	return "memory"
}

func (m *Memory) makeQueue() queue {
	if m.PoolNum <= 0 {
		return make(queue)
	}
	return make(queue, m.PoolNum)
}

// Publish 消息入生产者
func (m *Memory) Publish(ctx context.Context, message storage.Messager, optionFuncs ...func(*PublishOptions)) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	memoryMessage := new(Message)
	memoryMessage.SetId(message.GetId())
	memoryMessage.SetRoutingKey(message.GetRoutingKey())
	memoryMessage.SetValues(message.GetValues())
	v, ok := m.queue.Load(message.GetRoutingKey())

	// TODO: 错误超出5次将放弃
	if !ok && memoryMessage.GetErrorCount() < 5 {
		v = m.makeQueue()
		m.queue.Store(message.GetRoutingKey(), v)
		memoryMessage.SetErrorCount()
	}

	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(message.GetRoutingKey(), q)
	}
	go func(gm storage.Messager, gq queue) {
		gm.SetId(uuid.New().String())
		gq <- gm
	}(memoryMessage, q)
	return nil
}

// Consumer 监听消费者
func (m *Memory) Consumer(ctx context.Context, name string, f storage.ConsumerFunc, optionFuncs ...func(*ConsumeOptions)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.queue.Load(name)
	if !ok {
		v = m.makeQueue()
		m.queue.Store(name, v)
	}
	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(name, q)
	}
	go func(out queue, gf storage.ConsumerFunc) {
		var err error
		for message := range q {
			err = gf(ctx, message)
			if err != nil {
				out <- message
				err = nil
			}
		}
	}(q, f)
}

func (m *Memory) Run(ctx context.Context) {
	m.wait.Add(1)
	m.wait.Wait()
}

func (m *Memory) Shutdown(ctx context.Context) {
	m.wait.Done()
}
