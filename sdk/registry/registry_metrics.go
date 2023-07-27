package registry

import metrics "github.com/jxo-me/gf-metrics"

type MetricsRegistry struct {
	registry[*metrics.Monitor]
}
