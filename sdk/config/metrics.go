package config

type Metrics struct {
	Enable          bool      `json:"enable" yaml:"enable"`
	Path            string    `json:"path" yaml:"path"`
	SlowTime        int32     `json:"slowTime" yaml:"slowTime"`
	RequestDuration []float64 `json:"requestDuration" yaml:"requestDuration"`
}
