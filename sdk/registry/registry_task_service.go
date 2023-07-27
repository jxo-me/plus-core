package registry

import (
	"github.com/jxo-me/plus-core/core/task"
)

type TaskServiceRegistry struct {
	registry[task.TasksService]
}
