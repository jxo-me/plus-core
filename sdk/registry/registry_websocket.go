package registry

import (
	"github.com/jxo-me/plus-core/pkg/ws"
)

type WebSocketRegistry struct {
	registry[*ws.Instance]
}
