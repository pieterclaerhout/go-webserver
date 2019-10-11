package jobrunner

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/pieterclaerhout/go-log"
)

type Job struct {
	Name    string
	inner   cron.Job
	status  uint32
	Status  string
	Latency string
	running sync.Mutex
}

const UNNAMED = "(unnamed)"

func New(job cron.Job) *Job {
	name := reflect.TypeOf(job).Name()
	if name == "Func" {
		name = UNNAMED
	}
	return &Job{
		Name:  name,
		inner: job,
	}
}

func (j *Job) StatusUpdate() string {
	if atomic.LoadUint32(&j.status) > 0 {
		j.Status = "RUNNING"
		return j.Status
	}
	j.Status = "IDLE"
	return j.Status

}

func (j *Job) Run() {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			log.StackTrace(err.(error))
		}
	}()

	if !selfConcurrent {
		j.running.Lock()
		defer j.running.Unlock()
	}

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&j.status, 1)
	j.StatusUpdate()

	defer j.StatusUpdate()
	defer atomic.StoreUint32(&j.status, 0)

	j.inner.Run()

	end := time.Now()
	j.Latency = end.Sub(start).String()

}
