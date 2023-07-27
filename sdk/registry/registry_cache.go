package registry

import (
	cacheLib "github.com/jxo-me/plus-core/core/v2/cache"
)

type CacheRegistry struct {
	registry[cacheLib.ICache]
}
