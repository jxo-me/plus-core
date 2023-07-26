package queue

import (
	"context"
	"github.com/jxo-me/plus-core/core/message"
)

type ConsumerFunc func(ctx context.Context, msg message.IMessage) error

// ConsumeOptions are used to describe how a new consumer will be created.
type ConsumeOptions struct {
	// rabbitmq
	BindingRoutingKeys []string
	BindingExchange    *BindingExchangeOptions
	Concurrency        int
	QOSPrefetch        int
	ConsumerName       string
	ConsumerAutoAck    bool
	// rocketmq
	GroupName         string
	MaxReconsumeTimes int32
	AutoCommit        bool
}
