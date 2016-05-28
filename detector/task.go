package detector

import (
	"bytes"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"qiniupkg.com/x/log.v7"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/robfig/cron.v2"
)

// TaskArgs is a task query struct
type TaskArgs struct {
	ID       string        `json:"id,omitempty"`
	Name     string        `json:"name,omitempty"`
	URL      string        `json:"url,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
	Creator  string        `json:"creator,omitempty"`
}

// Task Cron
type Task struct {
	// Task options
	ID          bson.ObjectId `bson:"_id,omitempty"         json:"id"`
	Name        string        `bson:"name,omitempty"        json:"name"`
	Description string        `bson:"description,omitempty" json:"description"`
	Creator     string        `bson:"creator,omitempty"     json:"creator"`
	CreatedAt   time.Time     `bson:"created_at"            json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"            json:"updated_at"`
	Interval    time.Duration `bson:"interval,omitempty"    json:"interval"`
	Enabled     bool          `bson:"enabled,omitempty"     json:"enabled"`
	Status      string        `bson:"status,omitempty"      json:"status"`

	// HTTP Request args
	URL         string              `bson:"url,omitempty"          json:"url"`
	Protocol    string              `bson:"protocol,omitempty"     json:"protocol"`
	Method      string              `bson:"method,omitempty"       json:"method"`
	ContentType string              `bson:"content_type,omitempty" json:"content_type"`
	UserAgent   string              `bson:"user_agent,omitempty"   json:"user_agent"`
	Username    string              `bson:"username,omitempty"     json:"username"`
	Password    string              `bson:"password,omitempty"     json:"password"`
	Body        map[string][]string `bson:"body,omitempty"         json:"body"`

	// HTTP Response results
	LastStatusCode int           `bson:"last_status,omitempty"  json:"last_status"`
	TimeLatency    time.Duration `bson:"time_latency,omitempty" json:"time_latency"`

	// Job Info
	Job Job `json:"job"`
}

// Job cron exec
type Job struct {
	ID      cron.EntryID `json:"id"`
	status  uint32
	Status  string        `json:"status"`
	Latency time.Duration `json:"latency"`
	running sync.Mutex

	Next time.Time `json:"next"`
	Prev time.Time `json:"prev"`
}

// JobStatusUpdate Update task cron job status
func (t *Task) JobStatusUpdate() string {
	if atomic.LoadUint32(&t.Job.status) > 0 {
		t.Job.Status = "RUNNING"
	} else {
		t.Job.Status = "IDLE"
	}

	entry := MainCron.Entry(t.Job.ID)
	t.Job.Prev = entry.Prev
	t.Job.Next = entry.Next

	return t.Job.Status
}

// Run auto trigger in cron
func (t *Task) Run() {
	start := time.Now()

	// If the job panics, just print a stack trace.
	// Don't let the whole process die.
	defer func() {
		if err := recover(); err != nil {
			var buf bytes.Buffer
			logger := log.New(&buf, "Runner Log: ", log.Lshortfile)
			logger.Panic(err, "\n", string(debug.Stack()))
		}
	}()

	if !selfConcurrent {
		t.Job.running.Lock()
		defer t.Job.running.Unlock()
	}

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&t.Job.status, 1)
	t.JobStatusUpdate()

	defer atomic.StoreUint32(&t.Job.status, 0)
	defer t.JobStatusUpdate()

	// TODO

	t.Job.Latency = time.Since(start)
}

// Start task
func (t *Task) Start() {
	entry := MainCron.Entry(t.Job.ID)
	if entry.ID == 0 && t.Interval > 0 {
		t.Job.ID = MainCron.Schedule(cron.Every(t.Interval), t)
	}
	t.Job.Next = entry.Next
	t.Job.Prev = entry.Prev
}

// Stop task
func (t *Task) Stop() {
	if t.Job.ID > 0 {
		MainCron.Remove(t.Job.ID)
		t.Job.ID = 0
	}
}
