package examples

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/plus-core/pkg/v2/puller"
	"github.com/jxo-me/plus-core/pkg/v2/puller/config"
	"github.com/jxo-me/plus-core/pkg/v2/puller/mq"
	"time"
)

// StartPuller 示例：聚合平台服务中启动拉单任务
func StartPuller(ctx context.Context) error {
	var conf config.PullerConfig
	//_ = configloader.Load("puller", &conf)
	mqDsn := "amqp://user:password@rabbitmq:5672/"

	db := g.DB()
	//cursorRepo := cursor.NewMysqlCursorRepo(nil)

	mqProducer, err := mq.NewMQProducer(ctx, mqDsn)
	if err != nil {
		return err
	}

	pullerSvc := puller.NewPullerService(conf, db, mqProducer)

	for _, vendor := range conf.Vendors {
		go func(v string) {
			for {
				_ = pullerSvc.Scheduler.RunVendorPull(ctx, v)
				time.Sleep(time.Minute * 1)
			}
		}(vendor)
	}
	return nil
}
