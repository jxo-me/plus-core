package registry

import (
	"github.com/jxo-me/plus-core/core/v2/task"
)

type MemoryServiceRegistry struct {
	registry[task.MemoryService]
}
