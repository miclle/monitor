package jobrunner

import (
	"bytes"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/robfig/cron.v2"
)

// Job struct
type Job struct {
	JobID   cron.EntryID `json:"job_id"`
	Name    string       `json:"name"`
	inner   cron.Job
	status  uint32
	Status  string        `json:"status"`
	Latency time.Duration `json:"latency"`
	running sync.Mutex

	// Next time the job will run, or the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next time.Time `json:"next"`

	// Prev is the last time this job was run, or the zero time if never.
	Prev time.Time `json:"prev"`
}

// StatusUpdate return job status
func (job *Job) StatusUpdate() string {
	if atomic.LoadUint32(&job.status) > 0 {
		job.Status = "RUNNING"
		return job.Status
	}
	job.Status = "IDLE"
	return job.Status
}

// Run execute ...
func (job *Job) Run() {
	start := time.Now()
	// If the job panics, just print a stack trace.
	// Don't let the whole process die.
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			logger := log.New(&buf, "JobRunner Log: ", log.Lshortfile)
			logger.Panic(err, "\n", string(debug.Stack()))
		}
	}()

	if !selfConcurrent {
		job.running.Lock()
		defer job.running.Unlock()
	}

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&job.status, 1)
	job.StatusUpdate()

	defer atomic.StoreUint32(&job.status, 0)
	defer job.StatusUpdate()

	job.inner.Run()

	job.Latency = time.Since(start)
}
