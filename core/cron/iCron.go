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

type ICron interface {
	String() string
	GetCron() *gcron.Cron
	GetRawJobs() []Job
	GetWorkers() map[string]*gcron.Entry

	Start(context.Context)
	AddJobs(...Job) ICron
	GetJobs() []*gcron.Entry
	StartJob(context.Context, string) bool
	StopJob(context.Context, string) bool
	Stop(context.Context)
}
