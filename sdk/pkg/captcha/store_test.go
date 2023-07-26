package captcha

import (
	"fmt"
	cacheLib "github.com/jxo-me/plus-core/core/cache"
	"github.com/jxo-me/plus-core/sdk/cache/memory"
	"github.com/mojocn/base64Captcha"
	"math/rand"
	"testing"
	"time"
)

var _expiration = 6000

func getStore(_ *testing.T) cacheLib.ICache {
	return memory.NewMemory()
}

func TestSetGet(t *testing.T) {
	s := NewCacheStore(getStore(t), "", _expiration)
	id := "captcha id"
	d := "random-string"
	s.Set(id, d)
	d2 := s.Get(id, false)
	if d2 != d {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestGetClear(t *testing.T) {
	s := NewCacheStore(getStore(t), "", _expiration)
	id := "captcha id"
	d := "932839jfffjkdss"
	s.Set(id, d)
	d2 := s.Get(id, true)
	if d != d2 {
		t.Errorf("saved %v, getDigitsClear returned got %v", d, d2)
	}
	d2 = s.Get(id, false)
	if d2 != "" {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}

func BenchmarkSetCollect(b *testing.B) {
	store := memory.NewMemory()
	b.StopTimer()
	d := "fdskfew9832232r"
	s := NewCacheStore(store, "", -1)
	ids := make([]string, 1000)
	for i := range ids {
		ids[i] = fmt.Sprintf("%d", rand.Int63())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			s.Set(ids[j], d)
		}
	}
}
func TestStore_SetGoCollect(t *testing.T) {
	s := NewCacheStore(getStore(t), "", -1)
	for i := 0; i <= 100; i++ {
		s.Set(fmt.Sprint(i), fmt.Sprint(i))
	}
}

func TestStore_CollectNotExpire(t *testing.T) {
	s := NewCacheStore(getStore(t), "", 36000)
	for i := 0; i < 50; i++ {
		s.Set(fmt.Sprint(i), fmt.Sprint(i))
	}

	// let background goroutine to go
	time.Sleep(time.Second)

	if v := s.Get("0", false); v != "0" {
		fmt.Println("v:", v)
		t.Error("cache store get failed")
	}
}

func TestNewCacheStore(t *testing.T) {
	type args struct {
		store      cacheLib.ICache
		expiration int
	}
	tests := []struct {
		name string
		args args
		want base64Captcha.Store
	}{
		{"", args{getStore(t), 36000}, nil},
		{"", args{getStore(t), 180000}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCacheStore(tt.args.store, "", tt.args.expiration); got == nil {
				t.Errorf("NewMemoryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
