package memory

import (
	"context"
	"github.com/google/uuid"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/message"
	"sync"
)

type queue chan messageLib.IMessage

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
func (m *Memory) Publish(ctx context.Context, msg messageLib.IMessage, optionFuncs ...func(*queueLib.PublishOptions)) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	memoryMessage := new(message.Message)
	memoryMessage.SetId(msg.GetId())
	memoryMessage.SetRoutingKey(msg.GetRoutingKey())
	memoryMessage.SetValues(msg.GetValues())
	v, ok := m.queue.Load(msg.GetRoutingKey())

	// TODO: 错误超出5次将放弃
	if !ok && memoryMessage.GetErrorCount() < 5 {
		v = m.makeQueue()
		m.queue.Store(msg.GetRoutingKey(), v)
		memoryMessage.SetErrorIncr()
	}

	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(msg.GetRoutingKey(), q)
	}
	go func(gm messageLib.IMessage, gq queue) {
		gm.SetId(uuid.New().String())
		gq <- gm
	}(memoryMessage, q)
	return nil
}

// Consumer 监听消费者
func (m *Memory) Consumer(ctx context.Context, name string, f queueLib.ConsumerFunc, optionFuncs ...func(*queueLib.ConsumeOptions)) {
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
	go func(out queue, gf queueLib.ConsumerFunc) {
		var err error
		for iMessage := range q {
			err = gf(ctx, iMessage)
			if err != nil {
				out <- iMessage
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
