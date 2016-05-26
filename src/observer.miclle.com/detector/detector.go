package detector

import (
	"time"

	"gopkg.in/robfig/cron.v2"

	"qiniupkg.com/x/log.v7"
)

// DefaultJobPoolSize Default job poll size
const DefaultJobPoolSize = 10

var (
	// MainCron is a globle cron varibale
	MainCron *cron.Cron

	workPermits chan struct{}

	// selfConcurrent is a single job allowed to run concurrently with itself?
	selfConcurrent bool
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// Start MainCron
// Start(pool int, concurrent int) (10, 1)
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

	log.Infof("%s[Godeye Cron] %v Started... %s \n", magenta, time.Now(), reset)
}

// Stop MainCron
func Stop() {
	go StopNow()
}

// StopNow MainCron
func StopNow() {
	log.Infof("%s[Godeye Cron] %v Stoped... %s \n", magenta, time.Now(), reset)
	if MainCron != nil {
		MainCron.Stop()
	}
}

// CronEntries return MainCron entries
func CronEntries() []cron.Entry {
	return MainCron.Entries()
}

// CronTasks return MainCron tasks
func CronTasks() []*Task {
	ents := MainCron.Entries()

	tasks := make([]*Task, len(ents))

	for k, v := range ents {
		task := v.Job.(*Task)
		task.Job.ID = v.ID
		task.Job.Next = v.Next
		task.Job.Prev = v.Prev

		tasks[k] = task
	}

	return tasks
}
