package detector

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/miclle/observer/util/mgoutil"
)

// taskMgr struct
type taskMgr struct {
	coll *mgo.Collection
}

// List return a detector task instances array
func (mgr *taskMgr) List() (tasks []Task, err error) {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)

	tasks = make([]Task, 0, 0)
	err = coll.Find(bson.M{}).All(&tasks)
	return
}

// Create is create detector taks func
func (mgr *taskMgr) Create(task *Task) error {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)

	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now

	return coll.Insert(task)
}

// Find return a detector task instance
func (mgr *taskMgr) Find(args *TaskArgs) (task Task, err error) {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)

	query := bson.M{}

	if args.ID != "" {
		query["_id"] = bson.ObjectIdHex(args.ID)
	}

	if args.Name != "" {
		query["name"] = args.Name
	}

	err = coll.Find(query).One(&task)
	return
}

// Update is update task func
func (mgr *taskMgr) Update(task *Task) error {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)
	task.UpdatedAt = time.Now().UTC()
	return coll.UpdateId(task.ID, task)
}

// Delete is remove task by name
func (mgr *taskMgr) Delete(args *TaskArgs) error {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)

	query := bson.M{}

	if args.ID != "" {
		query["_id"] = bson.ObjectIdHex(args.ID)
	}

	if args.Name != "" {
		query["name"] = args.Name
	}

	return coll.Remove(query)
}

// DeleteAll is remove all task
func (mgr *taskMgr) DeleteAll() (err error) {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)
	_, err = coll.RemoveAll(bson.M{})
	return
}
