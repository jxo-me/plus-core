package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

type Item struct {
	ActiveT      int64  `json:"active_t"`
	AgentID      int64  `json:"agent_id"`
	Chips        int64  `json:"chips"`
	Chips2SQLT   int64  `json:"chips2sql_t"`
	ChipsChangeT int64  `json:"chips_change_t"`
	UserTag      string `json:"user_tag"`
}

func TestSet(t *testing.T) {
	opt := redis.Options{
		Addr:     "test.com:6379",
		Password: "aa123456",
	}
	client := redis.NewClient(&opt)
	sync := redsync.New(goredis.NewPool(client))
	item := Item{
		ActiveT:      1701245434,
		AgentID:      3540,
		Chips:        100000,
		Chips2SQLT:   1701244586,
		ChipsChangeT: 1701245440,
		UserTag:      "g2gCbQAAABgyMDIzMTEyOTA0MTA0MDI5NzY1MzA5NTRhBA==",
	}
	item1 := Item{AgentID: 3540}
	var tests = []struct {
		Table string
		Key   string
		Val   string
	}{
		{Table: "testBucket", Key: "123456", Val: gconv.String(item)},
		{Table: "testBucket", Key: "123457", Val: gconv.String(item1)},
	}
	for i, test := range tests {
		key := fmt.Sprintf("%s:%d", test.Table, i)
		r := NewRedis(test.Table, sync.NewMutex(key, redsync.WithRetryDelay(60*time.Second)), client, false)
		ctx := context.Background()
		set, err := r.Set(ctx, test.Key, test.Val)
		if err != nil {
			t.Error(`r.Set fail`, test.Table, test.Key, test.Val, err.Error())
		}
		t.Log("set result:", set)
		val, err := r.Get(ctx, test.Key)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}
		if test.Val != val {
			t.Error(`get result fail`, test.Val, val)
		}
		var res Item
		err = json.Unmarshal([]byte(val), &res)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}
		t.Log("json.Unmarshal result:", res.AgentID)
	}
}

func TestGet(t *testing.T) {
	opt := redis.Options{
		Addr:     "test.com:6379",
		Password: "aa123456",
	}
	client := redis.NewClient(&opt)
	sync := redsync.New(goredis.NewPool(client))
	item := Item{
		ActiveT:      1701245434,
		AgentID:      3540,
		Chips:        100000,
		Chips2SQLT:   1701244586,
		ChipsChangeT: 1701245440,
		UserTag:      "g2gCbQAAABgyMDIzMTEyOTA0MTA0MDI5NzY1MzA5NTRhBA==",
	}
	item1 := Item{AgentID: 3540}
	var tests = []struct {
		Table string
		Key   string
		Val   string
	}{
		{Table: "testBucket", Key: "123456", Val: gconv.String(item)},
		{Table: "testBucket", Key: "123457", Val: gconv.String(item1)},
	}
	for i, test := range tests {
		key := fmt.Sprintf("%s:%d", test.Table, i)
		r := NewRedis(test.Table, sync.NewMutex(key, redsync.WithRetryDelay(60*time.Second)), client, false)
		ctx := context.Background()
		val, err := r.Get(ctx, test.Key)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}
		if test.Val != val {
			t.Error(`get result fail`, test.Val, val)
		}
		var res Item
		err = json.Unmarshal([]byte(val), &res)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}
		t.Log("json.Unmarshal result:", res.AgentID)
	}
}

func TestHas(t *testing.T) {
	opt := redis.Options{
		Addr:     "test.com:6379",
		Password: "aa123456",
	}
	client := redis.NewClient(&opt)
	sync := redsync.New(goredis.NewPool(client))
	item := Item{
		ActiveT:      1701245434,
		AgentID:      3540,
		Chips:        100000,
		Chips2SQLT:   1701244586,
		ChipsChangeT: 1701245440,
		UserTag:      "g2gCbQAAABgyMDIzMTEyOTA0MTA0MDI5NzY1MzA5NTRhBA==",
	}
	item1 := Item{AgentID: 3540}
	var tests = []struct {
		Table string
		Key   string
		Val   string
	}{
		{Table: "testBucket", Key: "123456", Val: gconv.String(item)},
		{Table: "testBucket", Key: "123457", Val: gconv.String(item1)},
	}
	for i, test := range tests {
		key := fmt.Sprintf("%s:%d", test.Table, i)
		r := NewRedis(test.Table, sync.NewMutex(key, redsync.WithRetryDelay(60*time.Second)), client, false)
		ctx := context.Background()
		val, err := r.Has(ctx, test.Key)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}

		t.Log("r.Has result:", val)
	}
}

func TestDel(t *testing.T) {
	opt := redis.Options{
		Addr:     "test.com:6379",
		Password: "aa123456",
	}
	client := redis.NewClient(&opt)
	sync := redsync.New(goredis.NewPool(client))
	item := Item{
		ActiveT:      1701245434,
		AgentID:      3540,
		Chips:        100000,
		Chips2SQLT:   1701244586,
		ChipsChangeT: 1701245440,
		UserTag:      "g2gCbQAAABgyMDIzMTEyOTA0MTA0MDI5NzY1MzA5NTRhBA==",
	}
	item1 := Item{AgentID: 3540}
	var tests = []struct {
		Table string
		Key   string
		Val   string
	}{
		{Table: "testBucket", Key: "123456", Val: gconv.String(item)},
		{Table: "testBucket", Key: "123457", Val: gconv.String(item1)},
	}
	for i, test := range tests {
		key := fmt.Sprintf("%s:%d", test.Table, i)
		r := NewRedis(test.Table, sync.NewMutex(key, redsync.WithRetryDelay(60*time.Second)), client, false)
		ctx := context.Background()
		val, err := r.Del(ctx, test.Key)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}

		t.Log("r.Del result:", val)
	}
}

func TestLen(t *testing.T) {
	opt := redis.Options{
		Addr:     "test.com:6379",
		Password: "aa123456",
	}
	client := redis.NewClient(&opt)
	sync := redsync.New(goredis.NewPool(client))
	item := Item{
		ActiveT:      1701245434,
		AgentID:      3540,
		Chips:        100000,
		Chips2SQLT:   1701244586,
		ChipsChangeT: 1701245440,
		UserTag:      "g2gCbQAAABgyMDIzMTEyOTA0MTA0MDI5NzY1MzA5NTRhBA==",
	}
	item1 := Item{AgentID: 3540}
	var tests = []struct {
		Table string
		Key   string
		Val   string
	}{
		{Table: "testBucket", Key: "123456", Val: gconv.String(item)},
		{Table: "testBucket", Key: "123457", Val: gconv.String(item1)},
	}
	for i, test := range tests {
		key := fmt.Sprintf("%s:%d", test.Table, i)
		r := NewRedis(test.Table, sync.NewMutex(key, redsync.WithRetryDelay(60*time.Second)), client, false)
		ctx := context.Background()
		val, err := r.Len(ctx)
		if err != nil {
			t.Error(`r.Get fail`, test.Table, test.Key, test.Val, err.Error())
		}

		t.Log("r.Len result:", val)
	}
}
