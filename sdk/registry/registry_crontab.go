package registry

import "github.com/jxo-me/plus-core/core/v2/cron"

type CrontabRegistry struct {
	registry[cron.ICron]
}
