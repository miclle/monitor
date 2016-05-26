package mgr

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"observer.miclle.com/detector"
	"observer.miclle.com/util/mgoutil"
)

// TaskMgr struct
type TaskMgr struct {
	coll *mgo.Collection
}

// NewTaskMgr return a observer task manager
func NewTaskMgr(host, name, mode string) (*TaskMgr, error) {

	session := mgoutil.Open(&mgoutil.Config{
		Host: host,
		DB:   name,
		Mode: mode,
		Coll: "observer_tasks",
	})

	if err := session.Coll.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true}); err != nil {
		return nil, err
	}

	if err := session.Coll.EnsureIndex(mgo.Index{Key: []string{"url"}, Unique: true}); err != nil {
		return nil, err
	}
	return &TaskMgr{coll: session.Coll}, nil
}

// List return a detector task instances array
func (mgr *TaskMgr) List() (tasks []detector.Task, err error) {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)

	tasks = make([]detector.Task, 0, 0)
	err = coll.Find(bson.M{}).All(&tasks)
	return
}

// Create is create detector taks func
func (mgr *TaskMgr) Create(task *detector.Task) error {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)

	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now

	return coll.Insert(task)
}

// Find return a detector task instance
func (mgr *TaskMgr) Find(args *detector.TaskArgs) (task detector.Task, err error) {
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
func (mgr *TaskMgr) Update(task *detector.Task) error {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)
	task.UpdatedAt = time.Now().UTC()
	return coll.UpdateId(task.ID, task)
}

// Delete is remove task by name
func (mgr *TaskMgr) Delete(args *detector.TaskArgs) error {
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
func (mgr *TaskMgr) DeleteAll() (err error) {
	coll := mgoutil.FastCopyCollection(mgr.coll)
	defer mgoutil.CloseCollection(coll)
	_, err = coll.RemoveAll(bson.M{})
	return
}
