package cursor

import (
	"context"
	"time"
)

// CursorRepo 可落地到 MySQL/Redis/Mongo 实现
type CursorRepo interface {
	GetLastPullTime(ctx context.Context, vendor string) (time.Time, error)
	UpdatePullTime(ctx context.Context, vendor string, pullTime time.Time) error
}
