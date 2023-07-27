package registry

import (
	"github.com/jxo-me/plus-core/core/task"
)

type RocketMqServiceRegistry struct {
	registry[task.RocketMqService]
}
