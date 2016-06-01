package detector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskMgr(t *testing.T) {
	assert := assert.New(t)

	host := "127.0.0.1:27017"
	name := "observer_test"
	mode := "strong"

	var err error
	err = Init(host, name, mode)
	assert.Nil(err)
	assert.NotNil(TaskMgr)
}

func TestDeleteAll(t *testing.T) {
	assert := assert.New(t)
	err := TaskMgr.DeleteAll()
	assert.Nil(err)
}

func TestList(t *testing.T) {
	assert := assert.New(t)
	tasks, err := TaskMgr.List()
	assert.Nil(err)
	assert.NotNil(tasks)
	assert.True(len(tasks) == 0)
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	task := &Task{
		Name:        "Sample task name",
		Description: "Sample task description",
		URL:         "Sample task url",
		Protocol:    "Sample task protocol",
		Interval:    time.Minute,
	}

	err := TaskMgr.Create(task)
	assert.Nil(err)
	assert.NotNil(task)

	err = TaskMgr.Create(task)
	assert.NotNil(err)

	tasks, err := TaskMgr.List()
	assert.Nil(err)
	assert.NotNil(tasks)
	assert.True(len(tasks) > 0)
}

func TestFind(t *testing.T) {
	assert := assert.New(t)

	args := &TaskArgs{
		Name: "Sample task name",
	}

	task, err := TaskMgr.Find(args)
	assert.Nil(err)
	assert.NotNil(task)
	assert.NotEmpty(task.ID.Hex())

	args = &TaskArgs{
		Name: "Not exsit task name",
	}

	_, err = TaskMgr.Find(args)
	assert.NotNil(err)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	args := &TaskArgs{
		Name: "Sample task name",
	}

	task, err := TaskMgr.Find(args)
	assert.Nil(err)
	assert.NotNil(task)

	updateAt := task.UpdatedAt

	task.Description = "Sample task new description"
	err = TaskMgr.Update(&task)
	assert.Nil(err)
	assert.True(int64(task.UpdatedAt.Sub(updateAt)) > 0)
	assert.True(task.Description == "Sample task new description")
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	args := &TaskArgs{
		Name: "Sample task name",
	}

	err := TaskMgr.Delete(args)
	assert.Nil(err)
}
