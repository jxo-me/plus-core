package redis

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis redis模式
func NewRedis(client *redis.Client, options *redis.Options) (*Redis, error) {
	if client == nil {
		client = redis.NewClient(options)
	}
	r := &Redis{
		client: client,
	}
	err := r.connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis cache implement
type Redis struct {
	client *redis.Client
}

func (*Redis) String() string {
	return "redis"
}

// connect connect test
func (r *Redis) connect() error {
	var err error
	_, err = r.client.Ping(context.TODO()).Result()
	return err
}

// Get from a key
func (r *Redis) Get(ctx context.Context, key string) (*gvar.Var, error) {
	result, err := r.client.Get(context.TODO(), key).Result()
	if err != nil {
		return nil, err
	}
	return g.NewVar(result), nil
}

// Set value with key and expire time
func (r *Redis) Set(ctx context.Context, key string, val interface{}, expire int) error {
	return r.client.Set(context.TODO(), key, val, time.Duration(expire)*time.Second).Err()
}

// Del delete key in redis
func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(context.TODO(), key).Err()
}

// HashGet from a key
func (r *Redis) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	result, err := r.client.HGet(context.TODO(), hk, key).Result()
	if err != nil {
		return nil, err
	}
	return g.NewVar(result), nil
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(ctx context.Context, hk, key string) error {
	return r.client.HDel(context.TODO(), hk, key).Err()
}

func (r *Redis) Increase(ctx context.Context, key string) error {
	return r.client.Incr(context.TODO(), key).Err()
}

func (r *Redis) Decrease(ctx context.Context, key string) error {
	return r.client.Decr(context.TODO(), key).Err()
}

// Expire Set ttl
func (r *Redis) Expire(ctx context.Context, key string, dur time.Duration) error {
	return r.client.Expire(context.TODO(), key, dur).Err()
}

// GetClient 暴露原生client
func (r *Redis) GetClient() *redis.Client {
	return r.client
}
