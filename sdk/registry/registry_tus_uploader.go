package registry

import (
	"github.com/jxo-me/plus-core/pkg/v2/tus"
)

type TusRegistry struct {
	registry[*tus.Uploader]
}
