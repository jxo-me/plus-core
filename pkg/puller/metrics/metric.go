package metrics

import (
	"context"
	"github.com/gogf/gf/v2/os/gmetric"
)

func RecordWindowDuration(ctx context.Context, seconds float64) {
	counter, err := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{}).Counter("puller.window.duration", gmetric.MetricOption{})
	if err != nil {
	}
	counter.Add(ctx, seconds)
	//gmetric.GetDefault().Timer("puller.window.duration").Set(seconds)
}

func RecordWindowError(ctx context.Context) {
	counter, err := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{}).Counter("puller.window.error", gmetric.MetricOption{})
	if err != nil {
	}
	counter.Inc(ctx)
	//gmetric.GetDefault().Counter("puller.window.error").Inc()
}
