package jobrunner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type taskJob struct {
	Summary    string
	LastExecAt time.Time
}

func (j *taskJob) Name() string {
	return j.Summary
}

func (j *taskJob) Run() {
	j.LastExecAt = time.Now()
}

func TestRunner(t *testing.T) {

	assert := assert.New(t)
	assert.Nil(MainCron)
	assert.Nil(workPermits)
	assert.False(selfConcurrent)

	Start()
	assert.NotNil(MainCron)

	Start(10, 1)
	assert.NotNil(workPermits)
	assert.True(selfConcurrent)

	assert.Len(Entries(), 0)
	assert.Len(Jobs(), 0)

	id, err := Schedule("", nil)

	assert.Equal(id, 0)
	assert.NotNil(err)

	task := &taskJob{
		Summary: "test task",
	}

	id, err = Schedule("@every 1s", nil)

	assert.Equal(id, 0)
	assert.NotNil(err)

	id, err = Schedule("", task)

	assert.Equal(id, 0)
	assert.NotNil(err)

	id, err = Schedule("@every 1s", task)

	assert.Equal(id, 1)
	assert.Nil(err)
	assert.Len(Entries(), 1)
	assert.Len(Jobs(), 1)

	id = Every(time.Second, task)
	assert.Equal(id, 2)
	assert.Len(Entries(), 2)
	assert.Len(Jobs(), 2)

	Remove(id)
	assert.Len(Entries(), 1)
	assert.Len(Jobs(), 1)
}
