package model

import "time"

type CompensateMessage struct {
	Vendor string    `json:"vendor"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
