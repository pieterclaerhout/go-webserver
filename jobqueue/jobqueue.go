package jobqueue

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/jobqueue/jobrunner"
)

type JobQueue struct{}

type JobEntry interface {
	Run()
	Name() string
}

func Default() JobQueue {
	return JobQueue{}
}

func (jobQueue JobQueue) Start(poolSize int, concurrency int) {
	log.Info("Starting job queue, pool size:", poolSize, "concurrency:", concurrency)
	jobrunner.Start(poolSize, concurrency)
}

func (jobQueue JobQueue) Stop() {
	log.Info("Stopping job queue")
	jobrunner.Stop()
}

func (jobQueue JobQueue) Remove(id cron.EntryID) {
	jobrunner.Remove(id)
}

func (jobQueue JobQueue) Now(job JobEntry) {
	log.Info("Running job:", job.Name(), "immediately")
	jobrunner.Now(job)
}

func (jobQueue JobQueue) In(duration time.Duration, job JobEntry) {
	log.Info("Running job", job.Name(), "in:", duration)
	jobrunner.In(duration, job)
}

func (jobQueue JobQueue) Every(duration time.Duration, job JobEntry) {
	log.Info("Running job", job.Name(), "every:", duration)
	jobrunner.Every(duration, job)
}

func (jobQueue JobQueue) Schedule(spec string, job JobEntry) error {
	log.Info("Running job", job.Name(), "schedule:", spec)
	return jobrunner.Schedule(spec, job)
}

func (jobqueue JobQueue) StatusJSON() map[string]interface{} {
	return jobrunner.StatusJSON()
}
