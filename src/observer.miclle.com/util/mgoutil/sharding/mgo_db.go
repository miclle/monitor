package sharding

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"observer.miclle.com/util/mgoutil"
)

// Config sharding mgo configuration
type Config struct {
	Hosts          []string `json:"hosts"`
	DB             string   `json:"db"`
	Coll           string   `json:"coll"`
	Mode           string   `json:"mode"`
	Direct         int32    `json:"direct"`  // the value 0，then hosts
	SyncTimeoutInS int32    `json:"timeout"` // unit: second
}

// Session mgo session struct
type Session struct {
	*mgo.Session
	DB   *mgo.Database
	Coll *mgo.Collection
}

// Open return the database session
func Open(config *Config) *Session {

	info := &mgo.DialInfo{
		Addrs:   config.Hosts,
		Direct:  config.Direct != 0,
		Timeout: time.Duration(int64(config.SyncTimeoutInS) * 1e9),
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal("Connect MongoDB failed:", err, config.Hosts)
	}

	if config.Mode != "" {
		mgoutil.SetMode(session, config.Mode, true)
	}

	db := session.DB(config.DB)
	c := db.C(config.Coll)

	return &Session{session, db, c}
}
