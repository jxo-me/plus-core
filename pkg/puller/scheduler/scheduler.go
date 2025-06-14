package scheduler

import (
	"context"
	"time"

	"github.com/jxo-me/plus-core/pkg/v2/puller/cursor"
	"github.com/jxo-me/plus-core/pkg/v2/puller/dispatcher"
	"github.com/jxo-me/plus-core/pkg/v2/puller/worker"
)

type Scheduler struct {
	Dispatcher *dispatcher.Dispatcher
	CursorRepo cursor.CursorRepo
	Worker     *worker.PullWorker
}

func (s *Scheduler) RunVendorPull(ctx context.Context, vendor string) error {
	lastTime, err := s.CursorRepo.GetLastPullTime(ctx, vendor)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Add(-10 * time.Minute) // 留一定延迟
	windowSize := 15 * time.Minute

	for lastTime.Before(now) {
		nextEnd := lastTime.Add(windowSize)
		if nextEnd.After(now) {
			nextEnd = now
		}

		// 直接调用worker执行一次拉单窗口
		err = s.Worker.PullOneWindow(ctx, vendor, lastTime, nextEnd)
		if err != nil {
			return err
		}

		// 更新游标
		if err = s.CursorRepo.UpdatePullTime(ctx, vendor, nextEnd); err != nil {
			return err
		}

		lastTime = nextEnd
	}

	return nil
}
