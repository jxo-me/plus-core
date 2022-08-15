package captcha

import (
	"context"
	"github.com/jxo-me/plus-core/sdk/storage"
)

type cacheStore struct {
	cache      storage.AdapterCache
	expiration int
}

// NewCacheStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewCacheStore(cache storage.AdapterCache, expiration int) Store {
	s := new(cacheStore)
	s.cache = cache
	s.expiration = expiration
	return s
}

// Set sets the digits for the captcha id.
func (e *cacheStore) Set(ctx context.Context, id string, value string) error {
	return e.cache.Set(ctx, id, value, e.expiration)
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (e *cacheStore) Get(ctx context.Context, id string, clear bool) string {
	v, err := e.cache.Get(ctx, id)
	if err == nil {
		if clear {
			_ = e.cache.Del(ctx, id)
		}
		return v.String()
	}
	return ""
}

// Verify captcha's answer directly
func (e *cacheStore) Verify(ctx context.Context, id, answer string, clear bool) bool {
	return e.Get(ctx, id, clear) == answer
}
