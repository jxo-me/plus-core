package registry

import queueLib "github.com/jxo-me/plus-core/core/queue"

type QueueRegistry struct {
	registry[queueLib.IQueue]
}
