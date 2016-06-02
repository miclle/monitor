package detector

import (
	"bytes"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/satori/go.uuid"
	"qiniupkg.com/x/log.v7"
	"qiniupkg.com/x/xlog.v7"

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
	URL         string            `bson:"url,omitempty"          json:"url"`
	Protocol    string            `bson:"protocol,omitempty"     json:"protocol"`
	Method      string            `bson:"method,omitempty"       json:"method"`
	ContentType string            `bson:"content_type,omitempty" json:"content_type"`
	UserAgent   string            `bson:"user_agent,omitempty"   json:"user_agent"`
	Username    string            `bson:"username,omitempty"     json:"username"`
	Password    string            `bson:"password,omitempty"     json:"password"`
	Body        string            `bson:"body,omitempty"         json:"body"`
	BodyForm    map[string]string `bson:"body_form,omitempty"    json:"body_form"`

	// Last HTTP Response results
	Response TaskResponse `json:"response"`

	// Job Info
	Job Job `json:"job"`
}

// TaskResponse the task result
type TaskResponse struct {
	StatusCode  int           `bson:"status_code"  json:"status_code"`
	TimeLatency time.Duration `bson:"time_latency" json:"time_latency"`
	Message     string        `bson:"message"      json:"message"`
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

// Exec the task
func (t *Task) Exec() (resp *http.Response, err error) {
	client := &Client{
		Client:   &http.Client{Transport: http.DefaultTransport},
		Delivery: t.Notify,
	}

	l := xlog.New(uuid.NewV4().String())

	switch t.Method {
	case "POST", "PUT":
		if t.ContentType == "application/x-www-form-urlencoded" {
			data := make(map[string][]string)
			for k, v := range t.BodyForm {
				data[k] = []string{v}
			}
			msg := url.Values(data).Encode()
			resp, err = client.PostWith(l, t.URL, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))

		} else {
			resp, err = client.PostWith(l, t.URL, t.ContentType, strings.NewReader(t.Body), len(t.Body))
		}

	default:
		req, err := http.NewRequest(t.Method, t.URL, nil)
		if err != nil {
			return nil, err
		}
		resp, err = client.Do(l, req)
	}

	return nil, nil
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

// Notify messages
func (t *Task) Notify(l *xlog.Logger, start, end time.Time, req *http.Request, resp *http.Response, err error) {

	t.Response.TimeLatency = end.Sub(start)

	if resp != nil {
		t.Response.StatusCode = resp.StatusCode
	} else {
		t.Response.StatusCode = 0
	}

	if err != nil {
		t.Response.Message = err.Error()
	}

	// TODO: writing the response body to task response message field
	// note: need to limit in length

	// Update the Task to monitor status
	if err := TaskMgr.Update(t); err != nil {
		log.Error(err.Error())
	}

	// 4xx 5xx send a alarm
	if t.Response.StatusCode/100 == 4 || t.Response.StatusCode/100 == 5 {
		// TODO
	}
}
