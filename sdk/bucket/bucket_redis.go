package bucket

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"log"
)

type Redis struct {
	Mu    *redsync.Mutex
	Table string
	Store *redis.Client
}

func NewRedis(table string, mutex *redsync.Mutex, client *redis.Client) *Redis {
	return &Redis{
		Mu:    mutex,
		Table: table,
		Store: client,
	}
}

func (b *Redis) Index() string {
	return b.Table
}

func (b *Redis) Lock(ctx context.Context) error {
	return b.Mu.LockContext(ctx)
}

func (b *Redis) Unlock(ctx context.Context) {
	_, err := b.Mu.UnlockContext(ctx)
	if err != nil {
		log.Println("Redis UnlockContext error:", err.Error())
	}
}

func (b *Redis) Set(ctx context.Context, field string, value string) (int64, error) {
	err := b.Lock(ctx)
	if err != nil {
		return 0, err
	}
	defer b.Unlock(ctx)

	return b.Store.HSet(ctx, b.Table, field, value).Result()
}

func (b *Redis) Del(ctx context.Context, field string) (int64, error) {
	err := b.Lock(ctx)
	if err != nil {
		return 0, err
	}
	defer b.Unlock(ctx)

	return b.Store.HDel(ctx, b.Table, field).Result()
}

func (b *Redis) Get(ctx context.Context, field string) (string, error) {
	err := b.Lock(ctx)
	if err != nil {
		return "", err
	}
	defer b.Unlock(ctx)

	return b.Store.HGet(ctx, b.Table, field).Result()
}

func (b *Redis) Len(ctx context.Context) (int64, error) {
	err := b.Lock(ctx)
	if err != nil {
		return 0, err
	}
	defer b.Unlock(ctx)
	return b.Store.HLen(ctx, b.Table).Result()
}

func (b *Redis) Has(ctx context.Context, field string) (bool, error) {
	err := b.Lock(ctx)
	if err != nil {
		return false, err
	}
	defer b.Unlock(ctx)

	return b.Store.HExists(ctx, b.Table, field).Result()
}

func (b *Redis) All(ctx context.Context, field string) (map[string]string, error) {
	err := b.Lock(ctx)
	if err != nil {
		return nil, err
	}
	defer b.Unlock(ctx)

	return b.Store.HGetAll(ctx, field).Result()
}
