package puller

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/jxo-me/plus-core/pkg/v2/puller/adapter"
	"github.com/jxo-me/plus-core/pkg/v2/puller/config"
	"github.com/jxo-me/plus-core/pkg/v2/puller/cursor"
	"github.com/jxo-me/plus-core/pkg/v2/puller/dispatcher"
	"github.com/jxo-me/plus-core/pkg/v2/puller/mq"
	"github.com/jxo-me/plus-core/pkg/v2/puller/scheduler"
	"github.com/jxo-me/plus-core/pkg/v2/puller/worker"
)

type PullerService struct {
	Scheduler *scheduler.Scheduler
}

func NewPullerService(cfg config.PullerConfig, db gdb.DB, mqProducer mq.Producer) *PullerService {
	disp := dispatcher.NewDispatcher()

	// 厂商动态注册
	for _, vendor := range cfg.Vendors {
		switch vendor {
		case "TEST":
			disp.Register("TEST", adapter.NewTestAdapter())
		}

		cursorRepo := cursor.NewMysqlCursorRepo(db)
		pullWorker := worker.NewPullWorker(disp, mqProducer)

		return &PullerService{
			Scheduler: &scheduler.Scheduler{
				Dispatcher: disp,
				CursorRepo: cursorRepo,
				Worker:     pullWorker,
			},
		}
	}
	return nil
}
