package cron

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
)

type Job interface {
	GetSpec() *JobSpec
	Handle(ctx context.Context)
}

type JobSpec struct {
	Name    string
	Pattern string
}

type Adapter interface {
	String() string
	Start(context.Context)
	AddJobs(...Job) Adapter
	GetJobs() []*gcron.Entry
	StartJob(context.Context, string) bool
	StopJob(context.Context, string) bool
	Stop(context.Context)
}
