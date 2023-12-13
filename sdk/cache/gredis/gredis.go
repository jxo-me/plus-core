package gredis

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

// NewGredis redis模式
func NewGredis(client *gredis.Redis) (*Gredis, error) {
	r := &Gredis{
		client: client,
	}
	return r, nil
}

// Gredis cache implement
type Gredis struct {
	client *gredis.Redis
}

func (*Gredis) String() string {
	return "gredis"
}

// connect connect test
func (r *Gredis) connect() error {
	return nil
}

// Get from a key
func (r *Gredis) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return r.client.Do(ctx, "GET", key)
}

// Set value with key and expire time
func (r *Gredis) Set(ctx context.Context, key string, val interface{}, expire int) error {
	var err error
	if expire != 0 {
		_, err = r.client.Do(ctx, "SET", key, val, "EX", expire)
	} else {
		_, err = r.client.Do(ctx, "SET", key, val)
	}

	return err
}

// Del delete key in redis
func (r *Gredis) Del(ctx context.Context, key string) error {
	_, err := r.client.Do(ctx, "DEL", key)
	return err
}

// HashGet from a key
func (r *Gredis) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	return r.client.Do(ctx, "HGET", hk, key)
}

func (r *Gredis) HashSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error) {
	var s = []interface{}{key}
	for k, v := range fields {
		s = append(s, k, v)
	}
	v, err := r.client.Do(ctx, "HSet", s...)
	return v.Int64(), err
}

func (r *Gredis) HashMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	var s = []interface{}{key}
	for k, v := range fields {
		s = append(s, k, v)
	}
	_, err := r.client.Do(ctx, "HMSet", s...)
	return err
}

func (r *Gredis) HashMGet(ctx context.Context, key string, fields ...string) (gvar.Vars, error) {
	v, err := r.client.Do(ctx, "HMGet", append([]interface{}{key}, gconv.Interfaces(fields)...)...)
	return v.Vars(), err
}

func (r *Gredis) HashLen(ctx context.Context, key string) (int64, error) {
	v, err := r.client.Do(ctx, "HLen", key)
	return v.Int64(), err
}

func (r *Gredis) HashGetAll(ctx context.Context, key string) (*gvar.Var, error) {
	v, err := r.client.Do(ctx, "HGetAll", key)
	return v, err
}

// HashDel delete key in specify redis's hashtable
func (r *Gredis) HashDel(ctx context.Context, hk, key string) error {
	_, err := r.client.Do(ctx, "HDEL", hk, key)
	return err
}

// Increase get Increase
func (r *Gredis) Increase(ctx context.Context, key string) (int64, error) {
	v, err := r.client.Do(ctx, "INCR", key)
	return v.Int64(), err
}

func (r *Gredis) Decrease(ctx context.Context, key string) (int64, error) {
	v, err := r.client.Do(ctx, "DECR", key)
	return v.Int64(), err
}

// Expire Set ttl
func (r *Gredis) Expire(ctx context.Context, key string, dur time.Duration) error {
	_, err := r.client.Do(ctx, "EXPIRE", key, dur)
	return err
}

// GetClient 暴露原生client
func (r *Gredis) GetClient() *gredis.Redis {
	return r.client
}
