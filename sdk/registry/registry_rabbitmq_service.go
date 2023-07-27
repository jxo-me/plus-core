package registry

import (
	"github.com/jxo-me/plus-core/core/task"
)

type RabbitMqServiceRegistry struct {
	registry[task.RabbitMqService]
}
