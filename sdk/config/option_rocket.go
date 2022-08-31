package config

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type RocketOptions struct {
	Urls []string `yaml:"urls" json:"urls"`
	*primitive.Credentials
	LogPath   string `yaml:"logPath" json:"log_path"`
	LogFile   string `yaml:"logFile" json:"log_file"`
	LogLevel  string `yaml:"logLevel" json:"log_level"`
	LogStdout bool   `yaml:"logStdout" json:"log_stdout"`
}

func (e *RocketOptions) GetRocketOptions(ctx context.Context, s *Settings) (*RocketOptions, error) {
	opt, err := s.Cfg().Get(ctx, fmt.Sprintf("settings.queue.rocketmq.%s", DefaultGroupName), "")
	if err != nil {
		return nil, err
	}
	err = opt.Scan(&e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
