package worker

import (
	"context"
	"github.com/jxo-me/plus-core/pkg/v2/puller/dispatcher"
	"github.com/jxo-me/plus-core/pkg/v2/puller/logger"
	"github.com/jxo-me/plus-core/pkg/v2/puller/metrics"
	"github.com/jxo-me/plus-core/pkg/v2/puller/mq"
	"time"
)

type PullWorker struct {
	Dispatcher *dispatcher.Dispatcher
	Producer   mq.Producer
}

func NewPullWorker(dispatcher *dispatcher.Dispatcher, producer mq.Producer) *PullWorker {
	return &PullWorker{Dispatcher: dispatcher, Producer: producer}
}

func (w *PullWorker) PullOneWindow(ctx context.Context, vendor string, start, end time.Time) error {
	startTime := time.Now()
	records, err := w.Dispatcher.Pull(ctx, vendor, start, end)
	cost := time.Since(startTime).Seconds()

	metrics.RecordWindowDuration(ctx, cost)

	if err != nil {
		// 可加入异常通知、告警、MQ补单等逻辑
		metrics.RecordWindowError(ctx)
		_ = w.Producer.PushCompensateTask(ctx, vendor, start, end)
		logger.Errorf(ctx, "[失败补单] vendor:%s window:[%s~%s] err:%v", vendor, start.Format("15:04"), end.Format("15:04"), err)
		return err
	}

	// 记录日志与埋点
	logger.Infof(ctx, "[成功] vendor:%s window:[%s~%s] records:%d", vendor, start.Format("15:04"), end.Format("15:04"), len(records))
	// 后续可直接派发 downstream MQ 进行入库、清洗、统计

	return nil
}
