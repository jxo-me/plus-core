package queue

import (
	"context"
	"github.com/jxo-me/plus-core/core/message"
)

type IQueue interface {
	String() string
	Publish(ctx context.Context, message message.IMessage, optionFuncs ...func(*PublishOptions)) error
	Consumer(ctx context.Context, name string, f ConsumerFunc, optionFuncs ...func(*ConsumeOptions))
	Run(ctx context.Context)
	Shutdown(ctx context.Context)
}
