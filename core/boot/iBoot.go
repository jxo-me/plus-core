package boot

import (
	"context"
	"github.com/jxo-me/plus-core/core/v2/queue"
)

type Initialize interface {
	String() string
	Init(ctx context.Context) error
}

type QueueInitialize interface {
	Initialize
	GetQueue(ctx context.Context) (map[string]queue.IQueue, error)
}

type IBootstrap interface {
	Bootstrap(ctx context.Context, fs ...Initialize)
}
