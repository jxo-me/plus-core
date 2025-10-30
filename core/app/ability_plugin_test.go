package app

import (
	"sync"
	"testing"
)

func TestAbilityRegistry_RegisterAndGet(t *testing.T) {
	r := NewAbilityRegistry()
	r.Register("a", 1)
	if x, ok := r.Get("a"); !ok || x != 1 {
		t.Fatalf("get failed: %v", x)
	}
}

func TestAbilityRegistry_Remove(t *testing.T) {
	r := NewAbilityRegistry()
	r.Register("k", "v")
	r.Remove("k")
	if _, ok := r.Get("k"); ok {
		t.Fatal("expect removed")
	}
}

func TestAbilityRegistry_List(t *testing.T) {
	r := NewAbilityRegistry()
	for i := 0; i < 10; i++ {
		r.Register(string('a'+i), i)
	}
	names := r.List()
	if len(names) != 10 {
		t.Fatal("list length")
	}
}

func TestAbilityRegistry_Concurrent(t *testing.T) {
	r := NewAbilityRegistry()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Register(string('x'+(i%20)), i)
			_, _ = r.Get(string('x' + (i % 20)))
			r.List()
		}()
	}
	wg.Wait()
}
