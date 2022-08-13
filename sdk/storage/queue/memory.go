package queue

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/jxo-me/plus-core/sdk/storage"
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

func (m *Memory) Append(message storage.Messager) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	memoryMessage := new(Message)
	memoryMessage.SetID(message.GetID())
	memoryMessage.SetStream(message.GetStream())
	memoryMessage.SetValues(message.GetValues())

	v, ok := m.queue.Load(message.GetStream())

	// TODO: 错误超出5次将放弃
	if !ok && memoryMessage.GetErrorCount() < 5 {
		v = m.makeQueue()
		m.queue.Store(message.GetStream(), v)
		memoryMessage.SetErrorCount()
	}

	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(message.GetStream(), q)
	}
	go func(gm storage.Messager, gq queue) {
		gm.SetID(uuid.New().String())
		gq <- gm
	}(memoryMessage, q)
	return nil
}

func (m *Memory) Register(ctx context.Context, name string, f storage.ConsumerFunc) {
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

func (m *Memory) Run() {
	m.wait.Add(1)
	m.wait.Wait()
}

func (m *Memory) Shutdown() {
	m.wait.Done()
}
