package queue

import (
	"context"
	"fmt"
	"github.com/jxo-me/plus-core/core/message"
	queueLib "github.com/jxo-me/plus-core/core/queue"
	"github.com/jxo-me/plus-core/sdk/queue/memory"
	redis2 "github.com/jxo-me/plus-core/sdk/queue/redis"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/robinjoseph08/redisqueue/v2"
)

func TestNewMemoryQueue(t *testing.T) {
	type args struct {
		prefix string
		queue  queueLib.IQueue
	}
	q := memory.NewMemory(100)
	tests := []struct {
		name string
		args args
		want queueLib.IQueue
	}{
		{
			"test0",
			args{
				prefix: "",
				queue:  q,
			},
			&Queue{prefix: "", queue: q},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueue(tt.args.prefix, tt.args.queue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_Consumer(t *testing.T) {
	type fields struct {
		prefix string
		queue  queueLib.IQueue
	}
	type args struct {
		name string
		f    queueLib.ConsumerFunc
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "test.com:6379",
		Password: "aa123456",
	})
	q, err := redis2.NewRedis(&redisqueue.ProducerOptions{
		StreamMaxLength:      100,
		ApproximateMaxLength: true,
		RedisClient:          client,
	}, &redisqueue.ConsumerOptions{
		VisibilityTimeout: 60 * time.Second,
		BlockingTimeout:   5 * time.Second,
		ReclaimInterval:   1 * time.Second,
		BufferSize:        100,
		Concurrency:       10,
		RedisClient:       client,
	})
	if err != nil {
		t.Error(err)
		return
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test0",
			fields{
				prefix: "",
				queue:  q,
			},
			args{
				name: "operate_log_queue",
				f: func(ctx context.Context, m message.IMessage) error {
					fmt.Println(m.GetValues())
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//var e queueLib.IQueue
			e := &Queue{
				prefix: tt.fields.prefix,
				queue:  tt.fields.queue,
			}
			ctx := context.Background()
			e.Consumer(ctx, tt.args.name, tt.args.f)
			go e.Run(ctx)
		})
	}
}
