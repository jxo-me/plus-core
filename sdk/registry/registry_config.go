package registry

import (
	"github.com/gogf/gf/v2/os/gcfg"
)

type ConfigRegistry struct {
	registry[*gcfg.Config]
}
