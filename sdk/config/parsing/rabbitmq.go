package parsing

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/config"
	"github.com/jxo-me/plus-core/sdk/v2/queue/rabbitmq"
)

func ParseRabbitMq(ctx context.Context, cfg *config.RabbitmqOptions) (queueLib.IQueue, error) {
	logger := glog.New()
	err := logger.SetConfigWithMap(g.Map{
		"flags":  glog.F_TIME_STD | glog.F_FILE_LONG,
		"path":   cfg.LogPath,
		"file":   cfg.LogFile,
		"level":  cfg.LogLevel,
		"stdout": cfg.LogStdout,
	})

	q, err := rabbitmq.NewRabbitMQ(
		ctx,
		cfg.DSN,
		cfg.ReconnectInterval,
		cfg.Cfg,
		logger,
	)

	return q, err
}
