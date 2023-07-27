package registry

import (
	"github.com/jxo-me/plus-core/pkg/v2/ws"
)

type WebSocketRegistry struct {
	registry[*ws.Instance]
}
