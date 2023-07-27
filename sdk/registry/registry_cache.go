package registry

import (
	cacheLib "github.com/jxo-me/plus-core/core/cache"
)

type CacheRegistry struct {
	registry[cacheLib.ICache]
}
