package config

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/app"
	"github.com/jxo-me/plus-core/core/v2/boot"
)

const (
	QueueCfgName = "queueConfig"
)

var insQueue = Queue{}

type Queue struct {
	CfgList []boot.QueueInitialize
}

func QueueConfig() *Queue {
	return &insQueue
}

func (q *Queue) String() string {
	return QueueCfgName
}

func (q *Queue) Init(ctx context.Context, app app.IRuntime) error {
	s := app.ConfigRegister().Get(DefaultGroupName)
	rabbit, err := s.Get(ctx, "settings.queue.rabbitmq", "")
	if err != nil {
		return err
	}
	if rabbit.String() != "" {
		q.CfgList = append(q.CfgList, QueueRabbit())
	}
	memory, err := s.Get(ctx, "settings.queue.memory", "")
	if err != nil {
		return err
	}
	if memory.String() != "" {
		q.CfgList = append(q.CfgList, QueueMemory())
	}
	rocket, err := s.Get(ctx, "settings.queue.rocketmq", "")
	if err != nil {
		return err
	}
	if rocket.String() != "" {
		q.CfgList = append(q.CfgList, QueueRocket())
	}
	nsq, err := s.Get(ctx, "settings.queue.nsq", "")
	if err != nil {
		return err
	}
	if nsq.String() != "" {
		q.CfgList = append(q.CfgList, QueueNsq())
	}
	for _, queueCfg := range q.CfgList {
		err = queueCfg.Init(ctx, app)
		if err != nil {
			glog.Error(ctx, "Queue config init error:", err)
			return err
		}
	}

	return nil
}
