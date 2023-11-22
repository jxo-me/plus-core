package pkg

import "sync"

// Bucket 通用安全桶
type Bucket[T any] struct {
	mu  sync.Mutex
	Idx int
	M   map[string]T
}

func (b *Bucket[T]) Index() int {
	return b.Idx
}

func (b *Bucket[T]) Set(key string, item T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.M[key] = item
}

func (b *Bucket[T]) Del(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.M, key)
}

func (b *Bucket[T]) Get(key string) T {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.M[key]
}

func (b *Bucket[T]) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.M)
}

func (b *Bucket[T]) Has(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.M[key]; ok {
		return true
	}
	return false
}
