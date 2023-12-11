package send

import (
	"context"
	"fmt"
	"github.com/jxo-me/plus-core/core/v2/send"
)

const (
	Name = "MultiSend"
)

type MultiSend[T any] struct {
	Clients []send.ISender[T]
}

func NewMulti(c ...send.ISender[any]) *MultiSend[any] {
	return &MultiSend[any]{Clients: c}
}

func (t *MultiSend[T]) String() string {
	return Name
}

func (t *MultiSend[T]) Info(ctx context.Context, msg T) (err error) {
	for _, client := range t.Clients {
		err = client.Info(ctx, msg)
		if err != nil {
			return fmt.Errorf("%s send error:%s", client.String(), err)
		}
	}
	return nil
}

func (t *MultiSend[T]) Warn(ctx context.Context, msg T) (err error) {
	for _, client := range t.Clients {
		err = client.Warn(ctx, msg)
		if err != nil {
			return fmt.Errorf("%s send error:%s", client.String(), err)
		}
	}
	return nil
}

func (t *MultiSend[T]) Error(ctx context.Context, msg T) (err error) {
	for _, client := range t.Clients {
		err = client.Error(ctx, msg)
		if err != nil {
			return fmt.Errorf("%s send error:%s", client.String(), err)
		}
	}
	return nil
}
