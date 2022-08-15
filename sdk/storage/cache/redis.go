package cache

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"time"
)

// NewRedis redis模式
func NewRedis(client *gredis.Redis, options *gredis.Config) (*Redis, error) {
	var err error
	if client == nil {
		client, err = gredis.New(options)
		if err != nil {
			return nil, err
		}
	}
	r := &Redis{
		client: client,
	}
	return r, nil
}

// Redis cache implement
type Redis struct {
	client *gredis.Redis
}

func (*Redis) String() string {
	return "redis"
}

// connect connect test
func (r *Redis) connect() error {
	return nil
}

// Get from key
func (r *Redis) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return r.client.Do(ctx, "GET", key)
}

// Set value with key and expire time
func (r *Redis) Set(ctx context.Context, key string, val interface{}, expire int) error {
	_, err := r.client.Do(ctx, "SET", key, val, time.Duration(expire)*time.Second)
	return err
}

// Del delete key in redis
func (r *Redis) Del(ctx context.Context, key string) error {
	_, err := r.client.Do(ctx, "DEL", key)
	return err
}

// HashGet from key
func (r *Redis) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	return r.client.Do(ctx, "HGET", hk, key)
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(ctx context.Context, hk, key string) error {
	_, err := r.client.Do(ctx, "HDEL", hk, key)
	return err
}

// Increase get Increase
func (r *Redis) Increase(ctx context.Context, key string) error {
	_, err := r.client.Do(ctx, "INCRBY", key)
	return err
}

func (r *Redis) Decrease(ctx context.Context, key string) error {
	_, err := r.client.Do(ctx, "DECRBY", key)
	return err
}

// Expire Set ttl
func (r *Redis) Expire(ctx context.Context, key string, dur time.Duration) error {
	_, err := r.client.Do(ctx, "EXPIRE", key, dur)
	return err
}

// GetClient 暴露原生client
func (r *Redis) GetClient() *gredis.Redis {
	return r.client
}
