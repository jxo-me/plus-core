package memory

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemory_Get(t *testing.T) {
	type fields struct {
		items   *sync.Map
		queue   *sync.Map
		wait    sync.WaitGroup
		mutex   sync.RWMutex
		PoolNum uint
	}
	type args struct {
		key    string
		value  string
		expire int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"test01",
			fields{},
			args{
				key:    "test",
				value:  "test",
				expire: 10,
			},
			"test",
			false,
		},
		{
			"test02",
			fields{},
			args{
				key:    "test",
				value:  "test1",
				expire: 1,
			},
			"",
			true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory()
			if err := m.Set(ctx, tt.args.key, tt.args.value, tt.args.expire); err != nil {
				t.Errorf("Set() error = %v", err)
				return
			}
			time.Sleep(2 * time.Second)
			got, err := m.Get(ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
