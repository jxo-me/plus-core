package config

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type RocketOptions struct {
	Urls              []string `yaml:"urls" json:"urls"`
	GroupName         string   `yaml:"groupName" json:"group_name"`
	MaxReconsumeTimes int32    `yaml:"maxReconsumeTimes" json:"max_reconsume_times"`
	RetryTimes        int      `yaml:"retryTimes" json:"retry_times"`
	AccessKey         string   `yaml:"accessKey" json:"access_key"`
	SecretKey         string   `yaml:"secretKey" json:"secret_key"`
	Credentials       *primitive.Credentials
	AutoCommit        bool `yaml:"autoCommit" json:"auto_commit"`
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
