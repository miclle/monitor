// Package jobrunner is a job runner for executing scheduled or ad-hoc tasks asynchronously from HTTP requests.
//
// It adds a couple of features on top of the Robfig cron package:
// 1. Protection against job panics.  (They print to ERROR instead of take down the process)
// 2. (Optional) Limit on the number of jobs that may run simulatenously, to
//    limit resource consumption.
// 3. (Optional) Protection against multiple instances of a single job running
//    concurrently.  If one execution runs into the next, the next will be queued.
// 4. Cron expressions may be defined in app.conf and are reusable across jobs.
// 5. Job status reporting. [WIP]
package jobrunner

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/robfig/cron.v2"
)

// DefaultJobPoolSize default job poll size
const DefaultJobPoolSize = 10

var (
	// MainCron is singleton instance of the underlying job scheduler.
	MainCron *cron.Cron

	// This limits the number of jobs allowed to run concurrently.
	workPermits chan struct{}

	// Is a single job allowed to run concurrently with itself?
	selfConcurrent bool
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// Start the MainCron
// jobrunner.Start(pool int, concurrent int) (10, 1)
func Start(v ...int) {
	MainCron = cron.New()

	if len(v) > 0 {
		if v[0] > 0 {
			workPermits = make(chan struct{}, v[0])
		} else {
			workPermits = make(chan struct{}, DefaultJobPoolSize)
		}
	}

	if len(v) > 1 {
		if v[1] > 0 {
			selfConcurrent = true
		} else {
			selfConcurrent = false
		}
	}

	MainCron.Start()

	fmt.Printf("%s[JobRunner] %v Started... %s \n", magenta, time.Now().Format("2006/01/02 - 15:04:05"), reset)
}

// Stop ALL active jobs from running at the next scheduled time
func Stop() {
	go MainCron.Stop()
}

// Remove a specific job from running
// Get EntryID from the list job entries jobrunner.Entries()
// If job is in the middle of running, once the process is finished it will be removed
func Remove(id cron.EntryID) {
	MainCron.Remove(id)
}

// Entries return detailed list of currently running recurring jobs
// to remove an entry, first retrieve the ID of entry
func Entries() []cron.Entry {
	return MainCron.Entries()
}

// Jobs return StatusData array
func Jobs() []*Job {
	ents := MainCron.Entries()

	jobs := make([]*Job, len(ents))

	for k, v := range ents {
		job := v.Job.(*Job)
		job.JobID = v.ID
		job.Next = v.Next
		job.Prev = v.Prev
		jobs[k] = job
	}

	return jobs
}

// RunnerJob is an interface for submitted cron jobs.
type RunnerJob interface {
	Name() string
	Run()
}

// NewJob return job
func NewJob(job RunnerJob) *Job {
	return &Job{
		Name:  job.Name(),
		inner: job,
	}
}

// Schedule a cron job
func Schedule(spec string, j RunnerJob) (entryID cron.EntryID, err error) {
	sched, err := cron.Parse(spec)

	if err != nil {
		return 0, err
	}

	if j == nil {
		return 0, errors.New("Invalid memory address or nil pointer dereference")
	}

	job := NewJob(j)

	entryID = MainCron.Schedule(sched, job)

	job.JobID = entryID

	return
}

// Every run the given job at a fixed interval.
// The interval provided is the time between the job ending and the job being run again.
// The time that the job takes to run is not included in the interval.
func Every(duration time.Duration, j RunnerJob) cron.EntryID {
	job := NewJob(j)
	return MainCron.Schedule(cron.Every(duration), job)
}

// Now func run the given job right now.
func Now(job RunnerJob) {
	go NewJob(job).Run()
}

// In run the given job once, after the given delay.
func In(duration time.Duration, job RunnerJob) {
	go func() {
		time.Sleep(duration)
		NewJob(job).Run()
	}()
}
