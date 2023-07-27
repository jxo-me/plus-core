package registry

import (
	"github.com/jxo-me/plus-core/core/v2/task"
)

type RocketMqServiceRegistry struct {
	registry[task.RocketMqService]
}
