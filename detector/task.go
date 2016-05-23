package detector

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Task is a detector struct
type Task struct {
	ID          bson.ObjectId `bson:"_id,omitempty"          json:"id"`
	Name        string        `bson:"name,omitempty"         json:"name"`
	Description string        `bson:"description,omitempty"  json:"description"`
	URL         string        `bson:"url,omitempty"          json:"url"`
	Protocol    string        `bson:"protocol,omitempty"     json:"protocol"`
	Method      string        `bson:"method,omitempty"       json:"method"`
	Interval    time.Duration `bson:"interval,omitempty"     json:"interval"`
	LastStatus  int           `bson:"last_status,omitempty"  json:"last_status"`
	TimeLatency time.Duration `bson:"time_latency,omitempty" json:"time_latency"`
	CreatedAt   time.Time     `bson:"created_at"             json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"             json:"updated_at"`
}

// TaskArgs is a task query struct
type TaskArgs struct {
	ID       string        `json:"id,omitempty"`
	Name     string        `json:"name,omitempty"`
	URL      string        `json:"url,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
}
