package registry

import (
	"github.com/jxo-me/plus-core/core/task"
)

type MemoryServiceRegistry struct {
	registry[task.MemoryService]
}
