package registry

import "github.com/jxo-me/plus-core/core/v2/bucket"

type StateRegistry struct {
	registry[bucket.IState]
}
