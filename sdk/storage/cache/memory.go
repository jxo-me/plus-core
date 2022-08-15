package cache

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

// NewMemory memory模式
func NewMemory() *Memory {
	return &Memory{
		cache: gcache.New(),
	}
}

type Memory struct {
	cache *gcache.Cache
}

func (*Memory) String() string {
	return "memory"
}

func (m *Memory) connect() {
}

func (m *Memory) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return m.getItem(ctx, key)
}

func (m *Memory) getItem(ctx context.Context, key string) (*gvar.Var, error) {
	return m.cache.Get(ctx, key)
}

func (m *Memory) Set(ctx context.Context, key string, val interface{}, expire int) error {
	return m.setItem(ctx, key, val, expire)
}

func (m *Memory) setItem(ctx context.Context, key string, item interface{}, expire int) error {
	return m.cache.Set(ctx, key, item, time.Duration(expire)*time.Second)
}

func (m *Memory) Del(ctx context.Context, key string) error {
	return m.del(ctx, key)
}

func (m *Memory) del(ctx context.Context, key string) error {
	_, err := m.cache.Remove(ctx, key)
	return err
}

func (m *Memory) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	v, err := m.getItem(ctx, fmt.Sprintf("%s:%s", hk, key))
	if err != nil || v == nil {
		return nil, err
	}
	return v, err
}

func (m *Memory) HashDel(ctx context.Context, hk, key string) error {
	return m.del(ctx, fmt.Sprintf("%s:%s", hk, key))
}

func (m *Memory) Increase(ctx context.Context, key string) error {
	return m.calculate(ctx, key, 1)
}

func (m *Memory) Decrease(ctx context.Context, key string) error {
	return m.calculate(ctx, key, -1)
}

func (m *Memory) calculate(ctx context.Context, key string, num int) error {
	v, err := m.getItem(ctx, key)
	if err != nil {
		return err
	}
	if v == nil {
		return fmt.Errorf("%s not exist", key)
	}
	expire, err := m.cache.GetExpire(ctx, key)
	if err != nil {
		return err
	}
	return m.cache.Set(ctx, key, v.Int()+num, expire)
}

func (m *Memory) Expire(ctx context.Context, key string, dur time.Duration) error {
	v, err := m.getItem(ctx, key)
	if err != nil {
		return err
	}
	if v == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	return m.cache.Set(ctx, key, v, dur)
}
