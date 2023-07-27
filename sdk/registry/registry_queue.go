package registry

import queueLib "github.com/jxo-me/plus-core/core/v2/queue"

type QueueRegistry struct {
	registry[queueLib.IQueue]
}
