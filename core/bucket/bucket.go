package bucket

import "context"

type IBucketStore interface {
	Set(ctx context.Context, key string, value string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Len(ctx context.Context) (int64, error)
	Has(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) (int64, error)
	All(ctx context.Context) (map[string]string, error)
}

type IState interface {
	Set(ctx context.Context, id int64, item map[string]any) error
	Get(ctx context.Context, id int64) (map[string]any, error)
	Del(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, key string, val any) (bool, error)
	UpdateMap(ctx context.Context, id int64, data map[string]any) (map[string]any, error)
	Total(ctx context.Context) (int64, error)
	Has(ctx context.Context, id int64) (bool, error)
	All(ctx context.Context) (map[string]string, error)
	AllKeys(ctx context.Context) ([]int64, error)
}
