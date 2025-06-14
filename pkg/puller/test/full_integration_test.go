package test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/plus-core/pkg/v2/puller"
	"github.com/jxo-me/plus-core/pkg/v2/puller/config"
	"github.com/jxo-me/plus-core/pkg/v2/puller/mq"
	"testing"
)

func TestFullPullFlow(t *testing.T) {
	ctx := context.Background()
	mqDsn := "amqp://user:password@rabbitmq:5672/"

	db := g.DB()

	mqProducer, err := mq.NewMQProducer(ctx, mqDsn)
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.PullerConfig{
		Vendors:          []string{"TEST"},
		WindowSize:       15,
		Concurrency:      3,
		RetryCount:       1,
		TimeoutPerWindow: 20,
	}
	service := puller.NewPullerService(cfg, db, mqProducer)

	for _, vendor := range cfg.Vendors {
		if err = service.Scheduler.RunVendorPull(context.Background(), vendor); err != nil {
			fmt.Printf("Vendor %s 拉取失败: %v\n", vendor, err)
		}
	}
}
