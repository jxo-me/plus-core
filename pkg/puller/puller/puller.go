package puller

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/pkg/v2/puller/timeslider"
	"time"
)

// PullHandler 拉单回调函数签名
type PullHandler func(ctx context.Context, vendor string, windowStart, windowEnd time.Time) ([]interface{}, error)

// ExecutePullTask 拉单框架入口
func ExecutePullTask(ctx context.Context, cfg PullTaskConfig, handler PullHandler) ([]WindowPullResult, error) {
	sliderCfg := timeslider.SlidingWindowConfig{
		WindowSize:  cfg.WindowSize,
		Concurrency: cfg.Concurrency,
		RetryPolicy: timeslider.RetryPolicy{
			MaxRetries: cfg.RetryCount,
			RetryDelay: time.Second,
		},
		TimeoutPerWindow: cfg.TimeoutPerWindow,
	}

	var allResults []WindowPullResult

	results, err := timeslider.SlidingWindow(ctx, cfg.StartTime, cfg.EndTime, sliderCfg, func(ctx context.Context, windowStart, windowEnd time.Time) ([]interface{}, error) {
		data, err := handler(ctx, cfg.Vendor, windowStart, windowEnd)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, WindowPullResult{
			Vendor:      cfg.Vendor,
			WindowStart: windowStart,
			WindowEnd:   windowEnd,
			Records:     data,
		})
		return data, nil
	})

	if err != nil {
		return allResults, fmt.Errorf("pull failed: %w", err)
	}
	glog.Debug(ctx, "results:", results)
	return allResults, nil
}
