package jobrunner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type jobFunc struct {
	name string
	t    *testing.T
}

func (j *jobFunc) Run() {
	assert := assert.New(j.t)
	assert.Equal(j.name, "jobFuncName")
	j.name = "jobFuncNameModify"
}

func TestStatusUpdate(t *testing.T) {

	assert := assert.New(t)

	testJobFunc := &jobFunc{
		name: "jobFuncName",
		t:    t,
	}

	job := &Job{
		inner: testJobFunc,
	}

	status := job.StatusUpdate()
	assert.Equal(status, "IDLE")
	assert.True(job.status == 0)
	assert.Equal(job.Status, "IDLE")
	assert.Equal(job.Latency, 0)

	job.Run()

	assert.Equal(testJobFunc.name, "jobFuncNameModify")
	assert.Equal(status, "IDLE")
	assert.True(job.Latency > 0)
}
