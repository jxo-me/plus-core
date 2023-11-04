package queue

import (
	"context"
	"github.com/jxo-me/plus-core/core/v2/message"
)

type IQueue interface {
	String() string
	Publish(ctx context.Context, message message.IMessage, optionFuncs ...func(*PublishOptions)) error

	RpcRequest(ctx context.Context, key string, data []byte, optionFuncs ...func(*PublishOptions)) ([]byte, error)
	Consumer(ctx context.Context, name string, f ConsumerFunc, optionFuncs ...func(*ConsumeOptions))
	Run(ctx context.Context)
	Shutdown(ctx context.Context)
}
