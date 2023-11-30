package bucket

import "context"

type IBucketStore interface {
	Set(ctx context.Context, field string, value string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Len(ctx context.Context) (int64, error)
	Has(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) (int64, error)
}
