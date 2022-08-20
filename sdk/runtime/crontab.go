package runtime

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/sdk/cron"
)

type Crontab struct {
	Jobs    []cron.Job
	cron    *gcron.Cron
	Workers map[string]*gcron.Entry
}

func (t *Crontab) String() string {
	return "CrontabService"
}

func (t *Crontab) AddJobs(jobs ...cron.Job) {
	t.Jobs = jobs
}

func (t *Crontab) GetJobs() []*gcron.Entry {
	return t.cron.Entries()
}

func (t *Crontab) StartJob(ctx context.Context, jobName string) bool {
	if j, ok := t.Workers[jobName]; ok {
		j.Start()
		return true
	}
	return false
}

func (t *Crontab) StopJob(ctx context.Context, jobName string) bool {
	if j, ok := t.Workers[jobName]; ok {
		j.Stop()
		return true
	}
	return false
}

func (t *Crontab) Start(ctx context.Context) {
	var err error
	var entry *gcron.Entry
	t.cron = gcron.New()
	for _, job := range t.Jobs {
		sp := job.GetSpec()
		entry, err = t.cron.Add(ctx, sp.Pattern, job.Handle, sp.Name)
		if err != nil {
			glog.Debug(ctx, "cron job register error:", err.Error())
			continue
		}
		t.Workers[sp.Name] = entry
	}
}

func (t *Crontab) Stop(ctx context.Context) {
	for _, worker := range t.Workers {
		worker.Stop()
	}
}
