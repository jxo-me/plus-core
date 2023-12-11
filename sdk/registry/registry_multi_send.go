package registry

import (
	"github.com/jxo-me/plus-core/core/v2/send"
)

type SenderRegistry struct {
	registry[send.ISender[send.ISendMsg]]
}
