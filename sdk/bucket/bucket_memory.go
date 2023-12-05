package bucket

import (
	"context"
	"sync"
)

// Memory 通用安全桶
type Memory[T any] struct {
	Mu  sync.Locker
	Idx string
	M   map[string]T
}

func (b *Memory[T]) Index() string {
	return b.Idx
}

func (b *Memory[T]) Set(ctx context.Context, key string, item T) (int64, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	b.M[key] = item
	return 1, nil
}

func (b *Memory[T]) Del(ctx context.Context, key string) (int64, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	delete(b.M, key)
	return 1, nil
}

func (b *Memory[T]) Get(ctx context.Context, key string) (T, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	return b.M[key], nil
}

func (b *Memory[T]) Len(ctx context.Context) (int64, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	return int64(len(b.M)), nil
}

func (b *Memory[T]) Has(ctx context.Context, key string) (bool, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	if _, ok := b.M[key]; ok {
		return true, nil
	}
	return false, nil
}

func (b *Memory[T]) All(ctx context.Context, key string) (map[string]T, error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	return b.M, nil
}
