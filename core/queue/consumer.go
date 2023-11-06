package queue

import (
	"context"
	"github.com/jxo-me/plus-core/core/v2/message"
	"io"
)

type ConsumerFunc func(ctx context.Context, rw io.Writer, msg message.IMessage) error

// ConsumeOptions are used to describe how a new consumer will be created.
type ConsumeOptions struct {
	// rabbitmq
	BindingRoutingKeys []string
	BindingExchange    *BindingExchangeOptions
	Concurrency        int
	QOSPrefetch        int
	ConsumerName       string
	ConsumerAutoAck    bool
	Exclusive          bool

	// rocketmq
	GroupName         string
	MaxReconsumeTimes int32
	AutoCommit        bool
}
