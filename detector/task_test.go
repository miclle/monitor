package detector

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	assert := assert.New(t)

	StopCronNow()
	StartCron()

	err := TaskMgr.DeleteAll()
	assert.Nil(err)

	task := &Task{
		ID:          bson.NewObjectId(),
		Name:        "TestTask",
		Description: "This is test task description",
		Interval:    time.Second,
		URL:         "https://github.com",
		Method:      "GET",
	}

	TaskMgr.Create(task)

	assert.True(task.Job.ID == 0)

	task.Exec()

	assert.True(task.Response.TimeLatency > 0)
	assert.Equal(task.Response.StatusCode, 200)

	taskdb, err := TaskMgr.Find(&TaskArgs{
		ID: task.ID.Hex(),
	})

	assert.Nil(err)
	assert.NotNil(taskdb)
	assert.Equal(task.Response.TimeLatency, taskdb.Response.TimeLatency)
	assert.Equal(task.Response.StatusCode, taskdb.Response.StatusCode)
}
