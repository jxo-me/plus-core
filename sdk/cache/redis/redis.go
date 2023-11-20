package redis

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return g.NewVar(result), nil
}

// Set value with key and expire time
func (r *Redis) Set(ctx context.Context, key string, val interface{}, expire int) error {
	return r.client.Set(ctx, key, val, time.Duration(expire)*time.Second).Err()
}

// Del delete key in redis
func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// HashGet from a key
func (r *Redis) HashGet(ctx context.Context, hk, key string) (*gvar.Var, error) {
	result, err := r.client.HGet(ctx, hk, key).Result()
	if err != nil {
		return nil, err
	}
	return g.NewVar(result), nil
}

func (r *Redis) HashSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error) {
	var s []interface{}
	for k, v := range fields {
		s = append(s, k, v)
	}
	v, err := r.client.HSet(ctx, key, s...).Result()
	return v, err
}

func (r *Redis) HashMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	var s []interface{}
	for k, v := range fields {
		s = append(s, k, v)
	}
	_, err := r.client.HMSet(ctx, key, s...).Result()
	return err
}

func (r *Redis) HashMGet(ctx context.Context, key string, fields ...string) (gvar.Vars, error) {
	var vars gvar.Vars
	v, err := r.client.HMGet(ctx, key, fields...).Result()
	err = gconv.Structs(&v, &vars)
	if err != nil {
		return nil, err
	}
	return vars, err
}

func (r *Redis) HashLen(ctx context.Context, key string) (int64, error) {
	v, err := r.client.HLen(ctx, key).Result()
	return v, err
}

func (r *Redis) HashGetAll(ctx context.Context, key string) (*gvar.Var, error) {
	v, err := r.client.HGetAll(ctx, key).Result()
	return g.NewVar(v), err
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
