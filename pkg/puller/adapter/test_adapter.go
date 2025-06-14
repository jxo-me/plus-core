package adapter

import (
	"context"
	"fmt"
	"time"
)

type TestAdapter struct{}

func NewTestAdapter() *TestAdapter {
	return &TestAdapter{}
}

func (a *TestAdapter) Pull(ctx context.Context, startTime, endTime time.Time) ([]interface{}, error) {
	fmt.Printf("拉取 test [%s ~ %s]\n", startTime, endTime)
	// 实际对接你 Test SDK 实现
	// return testSdk.QueryBetRecordByTime(ctx, startTime, endTime)
	return []interface{}{"demo-test-record"}, nil
}
