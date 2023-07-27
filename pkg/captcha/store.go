package captcha

import (
	"context"
	"fmt"
	cacheLib "github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/mojocn/base64Captcha"
)

type cacheStore struct {
	cache      cacheLib.ICache
	expiration int
	prefix     string
}

// NewCacheStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewCacheStore(cache cacheLib.ICache, prefix string, expiration int) base64Captcha.Store {
	s := new(cacheStore)
	s.cache = cache
	s.prefix = prefix
	s.expiration = expiration
	return s
}

func (e *cacheStore) getPrefixKey(id string) string {
	return fmt.Sprintf("%s:%s", e.prefix, id)
}

// Set sets the digits for the captcha id.
func (e *cacheStore) Set(id string, value string) error {
	return e.cache.Set(context.Background(), e.getPrefixKey(id), value, e.expiration)
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (e *cacheStore) Get(id string, clear bool) string {
	v, err := e.cache.Get(context.Background(), e.getPrefixKey(id))
	if err == nil {
		if clear {
			_ = e.cache.Del(context.Background(), e.getPrefixKey(id))
		}
		return v.String()
	}
	return ""
}

// Verify captcha answer directly
func (e *cacheStore) Verify(id, answer string, clear bool) bool {
	return e.Get(id, clear) == answer
}
