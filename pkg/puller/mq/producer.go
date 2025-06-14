package mq

import (
	"context"
	"fmt"
	"time"
)

type Producer interface {
	PushCompensateTask(ctx context.Context, vendor string, start, end time.Time) error
}

type DummyProducer struct{}

func (d *DummyProducer) PushCompensateTask(ctx context.Context, vendor string, start, end time.Time) error {
	fmt.Printf("[补单派发] vendor=%s window=[%s ~ %s]\n", vendor, start, end)
	// 真实场景这里写入MQ
	return nil
}
