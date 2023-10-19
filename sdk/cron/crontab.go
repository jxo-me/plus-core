package cron

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/core/v2/cron"
)

var insCrontab = crontab{
	Jobs:    []cron.Job{},
	Cron:    gcron.New(),
	Workers: map[string]*gcron.Entry{},
}

type crontab struct {
	Jobs    []cron.Job
	Cron    *gcron.Cron
	Workers map[string]*gcron.Entry
}

func Crontab() *crontab {
	return &insCrontab
}

func (t *crontab) String() string {
	return "CrontabService"
}

func (t *crontab) GetCron() *gcron.Cron {
	return t.Cron
}

func (t *crontab) GetRawJobs() []cron.Job {
	return t.Jobs
}

func (t *crontab) GetWorkers() map[string]*gcron.Entry {
	return t.Workers
}

func (t *crontab) AddJobs(jobs ...cron.Job) cron.ICron {
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
		entry, err = t.Cron.AddSingleton(ctx, sp.Pattern, job.Handle, sp.Name)
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
