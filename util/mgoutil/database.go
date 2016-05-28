package mgoutil

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Database Usage
//     var database *mgo2.Database
//     func init() {
// 	    // when app start
// 	    var err error
// 	    database, err = mgo2.NewDatabase("mongodb://localhost/dbname", "strong")
// 	    if err != nil {
// 		    panic(err)
// 	    }
//     }
//     func foobar() {
// 	    db := database.Copy()
// 	    defer db.Close()
//     	// db.C(ColFoobar)...
//     }
type Database struct {
	session      *mgo.Session
	dbName       string
	indexService *MongoIndexService
}

// NewDatabase return database instance
func NewDatabase(url, mgoMode string) (res *Database, err error) {
	mgoAddr, dbName, err := parseMgoAddr(url)
	if err != nil {
		return
	}
	res, err = newDatabaseWithTimeout(mgoAddr, dbName, mgoMode, false)
	if err != nil {
		log.Fatal("<NewMongo> ", "mgo.Dial error:", err)
		return
	}
	return
}

// NewDatabaseWithTimeoutNoFatal return database instance
func NewDatabaseWithTimeoutNoFatal(url, mgoMode string, timeout time.Duration) (res *Database, err error) {
	mgoAddr, dbName, err := parseMgoAddr(url)
	if err != nil {
		return
	}
	return newDatabaseWithTimeout(mgoAddr, dbName, mgoMode, true, timeout)
}

func parseMgoAddr(url string) (mgoAddr string, dbName string, err error) {
	DbPos := strings.LastIndex(url, "/")
	if DbPos == -1 {
		err = errors.New("mgoDns don't contain '/'")
		return
	}
	mgoAddr = url[:DbPos]
	dbName = url[DbPos+1:]
	return
}

func newDatabaseWithTimeout(mgoAddr, dbName, mgoMode string, useTimeout bool, timeouts ...time.Duration) (res *Database, err error) {
	var mgoSession *mgo.Session
	var timeout time.Duration
	if useTimeout {
		if len(timeouts) > 0 {
			timeout = timeouts[0]
		}
		mgoSession, err = mgo.DialWithTimeout(mgoAddr, timeout)
	} else {
		mgoSession, err = mgo.Dial(mgoAddr)
	}
	if err != nil {
		err = fmt.Errorf("mgo.Dial error: %s", err)
		return
	}

	setMgoMode(mgoSession, mgoMode, true)

	res = &Database{
		session:      mgoSession,
		dbName:       dbName,
		indexService: NewMongoIndexService(mgoSession, dbName),
	}
	res.session.SetSyncTimeout(0)
	return
}

// C return the Collection
// for the format of col, check MongoIndexService.EnsureIndex
func (m *Database) C(col bson.M) *mgo.Collection {
	m.indexService.EnsureIndex(col)
	colName := col["name"].(string)
	return m.session.DB(m.dbName).C(colName)
}

// Close close the database session
func (m *Database) Close() {
	m.session.Close()
}

// Copy of the database
func (m *Database) Copy() *Database {
	return &Database{
		session:      m.session.Copy(),
		dbName:       m.dbName,
		indexService: m.indexService,
	}
}

// Session return the database session
func (m *Database) Session() *mgo.Session {
	return m.session
}
func setMgoMode(s *mgo.Session, modeFriendly string, refresh bool) {
	mode := strings.ToLower(modeFriendly)
	switch mode {
	case "eventual":
		s.SetMode(mgo.Eventual, refresh)
	case "monotonic", "mono":
		s.SetMode(mgo.Monotonic, refresh)
	case "strong":
		s.SetMode(mgo.Strong, refresh)
	default:
		log.Fatal("invalid mgo mode")
	}
}
