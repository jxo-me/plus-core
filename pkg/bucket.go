package pkg

import "sync"

// Bucket 通用安全桶
type Bucket[T any] struct {
	mu    sync.Mutex
	index int
	m     map[string]T
}

func (b *Bucket[T]) Index() int {
	return b.index
}

func (b *Bucket[T]) Add(key string, item T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.m[key] = item
}

func (b *Bucket[T]) Del(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.m, key)
}

func (b *Bucket[T]) Get(key string) T {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.m[key]
}

func (b *Bucket[T]) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.m)
}

func (b *Bucket[T]) Has(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.m[key]; ok {
		return true
	}
	return false
}
