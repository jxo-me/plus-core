package memory

import (
	"context"
	"fmt"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/message"
	"io"
	"log"
	"sync"
	"testing"
	"time"
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
		message messageLib.IMessage
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
				message: &message.Message{
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
		f    queueLib.ConsumerFunc
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
				f: func(ctx context.Context, rw io.Writer, msg messageLib.IMessage) error {
					fmt.Println(msg.GetValue())
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
			if err := m.Publish(ctx, &message.Message{
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
