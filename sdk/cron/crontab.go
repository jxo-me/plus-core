package cron

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
)

var insCrontab = crontab{
	Jobs:    []Job{},
	Cron:    gcron.New(),
	Workers: map[string]*gcron.Entry{},
}

type crontab struct {
	Jobs    []Job
	Cron    *gcron.Cron
	Workers map[string]*gcron.Entry
}

func Crontab() *crontab {
	return &insCrontab
}

func (t *crontab) String() string {
	return "CrontabService"
}

func (t *crontab) AddJobs(jobs ...Job) Adapter {
	t.Jobs = jobs
	return t
}

func (t *crontab) GetJobs() []*gcron.Entry {
	return t.Cron.Entries()
}

func (t *crontab) StartJob(ctx context.Context, jobName string) bool {
	if j, ok := t.Workers[jobName]; ok {
		j.Start()
		return true
	}
	return false
}

func (t *crontab) StopJob(ctx context.Context, jobName string) bool {
	if j, ok := t.Workers[jobName]; ok {
		j.Stop()
		return true
	}
	return false
}

func (t *crontab) Start(ctx context.Context) {
	var err error
	var entry *gcron.Entry
	for _, job := range t.Jobs {
		sp := job.GetSpec()
		entry, err = t.Cron.Add(ctx, sp.Pattern, job.Handle, sp.Name)
		if err != nil {
			glog.Debug(ctx, "cron job register error:", err.Error())
			continue
		}
		t.Workers[sp.Name] = entry
	}
}

func (t *crontab) Stop(ctx context.Context) {
	for _, worker := range t.Workers {
		worker.Stop()
	}
}
