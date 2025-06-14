package puller

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestExecutePullTask(t *testing.T) {
	ctx := context.Background()

	cfg := PullTaskConfig{
		Vendor:           "JILI",
		StartTime:        time.Date(2025, 6, 13, 0, 0, 0, 0, time.UTC),
		EndTime:          time.Date(2025, 6, 13, 1, 0, 0, 0, time.UTC),
		WindowSize:       15 * time.Minute,
		Concurrency:      3,
		RetryCount:       1,
		TimeoutPerWindow: 10 * time.Second,
	}

	mockHandler := func(ctx context.Context, vendor string, start, end time.Time) ([]interface{}, error) {
		fmt.Printf("Vendor: %s, Window: %s - %s\n", vendor, start.Format("15:04"), end.Format("15:04"))

		// 模拟部分窗口失败
		if start.Minute() == 30 {
			return nil, errors.New("simulate failed window")
		}
		return []interface{}{fmt.Sprintf("record-%s", start.Format("15:04"))}, nil
	}

	results, err := ExecutePullTask(ctx, cfg, mockHandler)

	// 打印输出
	for _, r := range results {
		t.Logf("Window [%s - %s] pulled %d records", r.WindowStart.Format("15:04"), r.WindowEnd.Format("15:04"), len(r.Records))
	}

	if err != nil {
		t.Logf("Partial errors captured: %v", err)
	}
}
