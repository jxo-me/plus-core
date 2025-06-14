package timeslider

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"

	"github.com/gogf/gf/v2/os/gmetric"
	"sync"
	"time"
)

// TimeWindowCallback 定义回调函数
type TimeWindowCallback func(ctx context.Context, windowStart, windowEnd time.Time) ([]interface{}, error)

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxRetries int           // 最大重试次数
	RetryDelay time.Duration // 每次重试间隔
}

// SlidingWindowConfig 滑动窗口配置
type SlidingWindowConfig struct {
	WindowSize       time.Duration // 每个时间窗口长度
	Concurrency      int           // 最大并发数
	RetryPolicy      RetryPolicy   // 重试策略
	TimeoutPerWindow time.Duration // 每个窗口最大处理时间
}

// DefaultSlidingWindowConfig 默认配置
func DefaultSlidingWindowConfig(windowSize time.Duration) SlidingWindowConfig {
	return SlidingWindowConfig{
		WindowSize:  windowSize,
		Concurrency: 5,
		RetryPolicy: RetryPolicy{
			MaxRetries: 2,
			RetryDelay: time.Second,
		},
		TimeoutPerWindow: 30 * time.Second,
	}
}

type taskResult struct {
	Result []interface{}
	Err    error
}

// SlidingWindow 核心滑动窗口执行器
func SlidingWindow(
	ctx context.Context,
	startTime, endTime time.Time,
	config SlidingWindowConfig,
	callback TimeWindowCallback,
) ([]interface{}, error) {

	if !startTime.Before(endTime) {
		return nil, errors.New("startTime must be before endTime")
	}
	if config.WindowSize <= 0 {
		return nil, errors.New("windowSize must be positive")
	}
	if config.Concurrency <= 0 {
		return nil, errors.New("concurrency must be positive")
	}

	taskCh := make(chan func() taskResult)
	resultCh := make(chan taskResult)

	// worker pool 启动
	var wg sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				resultCh <- task()
			}
		}()
	}

	// 任务生产
	go func() {
		for windowStart := startTime; windowStart.Before(endTime); windowStart = windowStart.Add(config.WindowSize) {
			windowEnd := windowStart.Add(config.WindowSize)
			if windowEnd.After(endTime) {
				windowEnd = endTime
			}
			start, end := windowStart, windowEnd
			taskCh <- func() taskResult {
				return executeWithRetry(ctx, start, end, config, callback)
			}
		}
		close(taskCh)
	}()

	// 结果收集
	var (
		results []interface{}
		errs    []error
		doneWg  sync.WaitGroup
	)
	doneWg.Add(1)
	go func() {
		defer doneWg.Done()
		for res := range resultCh {
			if res.Err != nil {
				errs = append(errs, res.Err)
			} else {
				results = append(results, res.Result...)
			}
		}
	}()

	wg.Wait()
	close(resultCh)
	doneWg.Wait()

	if len(errs) > 0 {
		return results, combineErrors(errs)
	}
	return results, nil
}

// 执行回调并内嵌重试机制
func executeWithRetry(
	ctx context.Context,
	windowStart, windowEnd time.Time,
	config SlidingWindowConfig,
	callback TimeWindowCallback,
) taskResult {

	spanCtx, span := gtrace.NewSpan(ctx, "SlidingWindowExecute")
	defer span.End()

	for attempt := 0; attempt <= config.RetryPolicy.MaxRetries; attempt++ {
		subCtx, cancel := context.WithTimeout(spanCtx, config.TimeoutPerWindow)
		start := time.Now()

		result, err := callback(subCtx, windowStart, windowEnd)
		duration := time.Since(start)

		cancel()

		//gmetric.GetDefault().Timer("timeslider.window.duration").Set(duration.Seconds())
		//if err == nil {
		//	return taskResult{Result: result, Err: nil}
		//}
		counter, err := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{}).Counter("timeslider.window.duration", gmetric.MetricOption{})
		if err != nil {
			return taskResult{Result: result, Err: nil}
		}
		counter.Add(ctx, duration.Seconds())
		g.Log().Warningf(ctx, "SlidingWindow attempt=%d failed window[%s~%s]: %v", attempt+1, windowStart.Format(time.RFC3339), windowEnd.Format(time.RFC3339), err)
		span.RecordError(err)

		if attempt < config.RetryPolicy.MaxRetries {
			time.Sleep(config.RetryPolicy.RetryDelay)
		}
	}

	return taskResult{Result: nil, Err: fmt.Errorf("failed after retries window[%s~%s]", windowStart, windowEnd)}
}

// 聚合错误信息
func combineErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	}
	errMsg := "multiple errors:"
	for _, e := range errs {
		errMsg += " [" + e.Error() + "]"
	}
	return errors.New(errMsg)
}
