package model

import "time"

type WindowTask struct {
	Vendor string
	Start  time.Time
	End    time.Time
}
