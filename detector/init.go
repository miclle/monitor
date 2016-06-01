package detector

import (
	"github.com/miclle/observer/util/mgoutil"
	"gopkg.in/mgo.v2"
)

// TaskMgr is a globle task manager
var TaskMgr *taskMgr

// Init return a observer task manager
func Init(host, name, mode string) (err error) {

	session := mgoutil.Open(&mgoutil.Config{
		Host: host,
		DB:   name,
		Mode: mode,
		Coll: "observer_tasks",
	})

	err = session.Coll.EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})

	TaskMgr = &taskMgr{coll: session.Coll}
	return
}
