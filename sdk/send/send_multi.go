package send

import (
	"context"
	"fmt"
	"github.com/jxo-me/plus-core/core/v2/send"
)

const (
	Name = "MultiSend"
)

type MultiSend[T send.ISendMsg] struct {
	clients []send.ISender[T]
}

func NewMulti(c ...send.ISender[send.ISendMsg]) *MultiSend[send.ISendMsg] {
	return &MultiSend[send.ISendMsg]{clients: c}
}

func (t *MultiSend[T]) String() string {
	return Name
}

func (t *MultiSend[T]) Info(ctx context.Context, msg T) (err error) {
	for _, client := range t.clients {
		err = client.Info(ctx, msg)
		if err != nil {
			return fmt.Errorf("%s send error:%s", client.String(), err)
		}
	}
	return nil
}

func (t *MultiSend[T]) Warn(ctx context.Context, msg T) (err error) {
	for _, client := range t.clients {
		err = client.Warn(ctx, msg)
		if err != nil {
			return fmt.Errorf("%s send error:%s", client.String(), err)
		}
	}
	return nil
}

func (t *MultiSend[T]) Error(ctx context.Context, msg T) (err error) {
	for _, client := range t.clients {
		err = client.Error(ctx, msg)
		if err != nil {
			return fmt.Errorf("%s send error:%s", client.String(), err)
		}
	}
	return nil
}
