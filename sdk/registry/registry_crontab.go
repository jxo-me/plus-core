package registry

import "github.com/jxo-me/plus-core/core/cron"

type CrontabRegistry struct {
	registry[cron.ICron]
}
