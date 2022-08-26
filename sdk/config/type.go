package config

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
)

type Initialize interface {
	String() string
	Init(ctx context.Context, s *Settings) error
}

type QueueInitialize interface {
	Initialize
	GetQueue(ctx context.Context) (storage.AdapterQueue, error)
}
