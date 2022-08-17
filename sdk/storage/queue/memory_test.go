package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/jxo-me/plus-core/sdk/storage"
)

func TestMemory_Publish(t *testing.T) {
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
	}
	type args struct {
		name    string
		message storage.Messager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test01",
			fields{},
			args{
				name: "test",
				message: &Message{
					Id:         "",
					RoutingKey: "test",
					Values: map[string]interface{}{
						"key": "value",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m := NewMemory(100)
			if err := m.Publish(ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Consumer(t *testing.T) {
	log.SetFlags(19)
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
	}
	type args struct {
		name string
		f    storage.ConsumerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test01",
			fields{},
			args{
				name: "test",
				f: func(ctx context.Context, message storage.Messager) error {
					fmt.Println(message.GetValues())
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(100)
			ctx := context.Background()
			m.Consumer(ctx, tt.name, tt.args.f)
			if err := m.Publish(ctx, &Message{
				Id:         "",
				RoutingKey: "test",
				Values: map[string]interface{}{
					"key": "value",
				}}); err != nil {
				t.Error(err)
				return
			}
			go func() {
				m.Run(ctx)
			}()
			time.Sleep(3 * time.Second)
			m.Shutdown(ctx)
		})
	}
}
