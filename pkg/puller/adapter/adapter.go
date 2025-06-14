package adapter

import (
	"context"
	"time"
)

type VendorPullAdapter interface {
	Pull(ctx context.Context, startTime, endTime time.Time) ([]interface{}, error)
}
