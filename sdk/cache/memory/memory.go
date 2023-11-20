package memory

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"strings"
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

func (m *Memory) HashSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error) {
	var n int64
	var err error
	for k, v := range fields {
		err = m.setItem(ctx, fmt.Sprintf("%s:%s", key, k), v, -1)
		if err != nil {
			continue
		}
		n++
	}

	return n, err
}

func (m *Memory) HashMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	var n int64
	var err error
	for k, v := range fields {
		err = m.setItem(ctx, fmt.Sprintf("%s:%s", key, k), v, -1)
		if err != nil {
			continue
		}
		n++
	}
	return err
}

func (m *Memory) HashMGet(ctx context.Context, key string, fields ...string) (gvar.Vars, error) {
	var vars gvar.Vars
	for _, hk := range fields {
		v, err := m.getItem(ctx, fmt.Sprintf("%s:%s", hk, key))
		if err != nil || v == nil {
			return nil, err
		}
		vars = append(vars, v)
	}

	return vars, nil
}

func (m *Memory) HashLen(ctx context.Context, key string) (int64, error) {
	var n int64
	keys, err := m.cache.Keys(ctx)
	if err != nil {
		return 0, err
	}
	for _, mkey := range keys {
		if s, ok := mkey.(string); ok {
			if strings.Contains(s, fmt.Sprintf("%s:", key)) {
				n++
			}
		}
	}
	return n, err
}

func (m *Memory) HashGetAll(ctx context.Context, key string) (*gvar.Var, error) {
	var vals map[string]interface{}
	keys, err := m.cache.Keys(ctx)
	if err != nil {
		return nil, err
	}
	for _, mkey := range keys {
		if s, ok := mkey.(string); ok {
			if strings.Contains(s, fmt.Sprintf("%s:", key)) {
				get, err := m.Get(ctx, s)
				if err != nil {
					return nil, err
				}
				list := strings.Split(s, ":")
				if len(list) > 0 {
					vals[list[1]] = get.Val()
				}
			}
		}
	}
	return g.NewVar(vals), err
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
