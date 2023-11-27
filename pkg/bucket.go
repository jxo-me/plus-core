package pkg

import "sync"

// Bucket 通用安全桶
type Bucket[T any] struct {
	Mu  sync.Locker
	Idx int
	M   map[string]T
}

func (b *Bucket[T]) Index() int {
	return b.Idx
}

func (b *Bucket[T]) Set(key string, item T) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	b.M[key] = item
}

func (b *Bucket[T]) Del(key string) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	delete(b.M, key)
}

func (b *Bucket[T]) Get(key string) T {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	return b.M[key]
}

func (b *Bucket[T]) Len() int {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	return len(b.M)
}

func (b *Bucket[T]) Has(key string) bool {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	if _, ok := b.M[key]; ok {
		return true
	}
	return false
}
