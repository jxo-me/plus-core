package registry

import (
	"github.com/jxo-me/plus-core/pkg/tus"
)

type TusRegistry struct {
	registry[*tus.Uploader]
}
