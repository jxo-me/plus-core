package config

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/queue/rocketmq"
)

const (
	RocketQueueName = "rocketmq"
)

type RocketOptions struct {
	Urls []string `yaml:"urls" json:"urls"`
	*primitive.Credentials
	LogPath   string `yaml:"logPath" json:"log_path"`
	LogFile   string `yaml:"logFile" json:"log_file"`
	LogLevel  string `yaml:"logLevel" json:"log_level"`
	LogStdout bool   `yaml:"logStdout" json:"log_stdout"`
}

var insQueueRocket = cQueueRocket{
	List: map[string]*RocketOptions{},
}

type cQueueRocket struct {
	List map[string]*RocketOptions
}

func QueueRocket() *cQueueRocket {
	return &insQueueRocket
}

func (c *cQueueRocket) String() string {
	return RocketQueueName
}

func (c *cQueueRocket) Init(ctx context.Context) error {
	var err error
	conf, err := Setting().Cfg().Get(ctx, "settings.queue.rabbitmq", "")
	if err != nil {
		return err
	}
	list := make(map[string]*RocketOptions)
	err = conf.Scan(&list)
	if err != nil {
		return err
	}
	for key, config := range list {
		// primitive.Credentials
		if config.AccessKey != "" && config.SecretKey != "" {
			config.Credentials = &primitive.Credentials{
				AccessKey: config.AccessKey,
				SecretKey: config.SecretKey,
			}
		}
		c.List[key] = config
	}

	return nil
}

// GetQueue get Rocket queue
func (c *cQueueRocket) GetQueue(ctx context.Context) (map[string]queueLib.IQueue, error) {
	list := make(map[string]queueLib.IQueue)
	for key, options := range c.List {
		logger := glog.New()
		err := logger.SetConfigWithMap(g.Map{
			"flags":  glog.F_TIME_STD | glog.F_FILE_LONG,
			"path":   options.LogPath,
			"file":   options.LogFile,
			"level":  options.LogLevel,
			"stdout": options.LogStdout,
		})
		if err != nil {
			return nil, err
		}
		list[key], err = rocketmq.NewRocketMQ(
			ctx,
			options.Urls,
			options.Credentials,
			logger,
		)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
