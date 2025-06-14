package dispatcher

import (
	"context"
	"errors"
	"github.com/jxo-me/plus-core/pkg/v2/puller/adapter"
	"time"
)

type Dispatcher struct {
	adapters map[string]adapter.VendorPullAdapter
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{adapters: make(map[string]adapter.VendorPullAdapter)}
}

func (d *Dispatcher) Register(vendor string, adapter adapter.VendorPullAdapter) {
	d.adapters[vendor] = adapter
}

func (d *Dispatcher) GetAdapter(vendor string) (adapter.VendorPullAdapter, error) {
	if a, ok := d.adapters[vendor]; ok {
		return a, nil
	}
	return nil, errors.New("vendor adapter not found")
}

func (d *Dispatcher) Pull(ctx context.Context, vendor string, start, end time.Time) ([]interface{}, error) {
	if a, ok := d.adapters[vendor]; ok {
		return a.Pull(ctx, start, end)
	}
	return nil, errors.New("vendor adapter not registered: " + vendor)
}
