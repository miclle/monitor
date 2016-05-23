package mgoutil

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

// CopySessionMaxRetry retry copy session max 5 times
var CopySessionMaxRetry = 5

// Dail return mgo session
func Dail(host string, mode string, syncTimeoutInS int64) *mgo.Session {

	session, err := mgo.Dial(host)
	if err != nil {
		log.Fatal("Connect MongoDB failed:", err, host)
	}

	if mode != "" {
		SetMode(session, mode, true)
	}
	if syncTimeoutInS != 0 {
		session.SetSyncTimeout(time.Duration(int64(time.Second) * syncTimeoutInS))
	}

	return session
}

// Safe mgo configuration
type Safe struct {
	W        int    `json:"w"`
	WMode    string `json:"wmode"`
	WTimeout int    `json:"wtimeoutms"`
	FSync    bool   `json:"fsync"`
	J        bool   `json:"j"`
}

// Config mgo configuration
type Config struct {
	Host           string `json:"host"`
	DB             string `json:"db"`
	Coll           string `json:"coll"`
	Mode           string `json:"mode"`
	SyncTimeoutInS int64  `json:"timeout"` // unit: second
	Safe           *Safe  `json:"safe"`
}

// Session mgo session struct
type Session struct {
	*mgo.Session
	DB   *mgo.Database
	Coll *mgo.Collection
}

// Open a mgo configuration return a session
func Open(cfg *Config) *Session {

	session := Dail(cfg.Host, cfg.Mode, cfg.SyncTimeoutInS)
	EnsureSafe(session, cfg.Safe)
	db := session.DB(cfg.DB)
	c := db.C(cfg.Coll)

	return &Session{session, db, c}
}

// IsSessionClosed test whether session closed
// PS: sometimes it's not corrected
func IsSessionClosed(s *mgo.Session) (res bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("[MGO_IS_SESSION_CLOSED] check session closed panic:", err)
		}
	}()
	res = true
	return s.Ping() != nil
}

func checkSession(s *mgo.Session) (err error) {
	return s.Ping()
}

func isServersFailed(err error) bool {
	return strings.Contains(err.Error(), "no reachable servers")
}

// CopySession return a copy mgo session
func CopySession(s *mgo.Session) *mgo.Session {
	for i := 0; i < CopySessionMaxRetry; i++ {
		res := s.Copy()
		err := checkSession(res)
		if err == nil {
			return res
		}
		CloseSession(res)
		log.Print("[MGO_COPY_SESSION] copy session and check failed:", err)
		if isServersFailed(err) {
			panic("[MGO_COPY_SESSION_FAILED] servers failed")
		}
	}
	msg := fmt.Sprintf("[MGO_COPY_SESSION_FAILED] failed after %d retries", CopySessionMaxRetry)
	log.Print(msg)
	panic(msg)
}

// FastCopySession return mgo session
func FastCopySession(s *mgo.Session) *mgo.Session {
	return s.Copy()
}

// CloseSession close the mgo session
func CloseSession(s *mgo.Session) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("[MGO_CLOSE_SESSION_RECOVER] close session panic", err)
		}
	}()
	s.Close()
}

// CopyDatabase copy database's session, and re-create the database.
// you need call `CloseDatabase` after use this
func CopyDatabase(db *mgo.Database) *mgo.Database {
	return CopySession(db.Session).DB(db.Name)
}

// FastCopyDatabase copy database's session
func FastCopyDatabase(db *mgo.Database) *mgo.Database {
	return FastCopySession(db.Session).DB(db.Name)
}

// CloseDatbase close the session of the database
func CloseDatbase(db *mgo.Database) {
	CloseSession(db.Session)
}

// CopyCollection copy collection's session, and re-create the collection
// you need call `CloseColletion` after use this
func CopyCollection(c *mgo.Collection) *mgo.Collection {
	return CopyDatabase(c.Database).C(c.Name)
}

// FastCopyCollection copy collection's session
func FastCopyCollection(c *mgo.Collection) *mgo.Collection {
	return FastCopyDatabase(c.Database).C(c.Name)
}

// CloseCollection close the session of the collection
func CloseCollection(c *mgo.Collection) {
	CloseDatbase(c.Database)
}

// CheckIndex check mgo collection index
func CheckIndex(c *mgo.Collection, key []string, unique bool) error {
	originIndexs, err := c.Indexes()
	if err != nil {
		return fmt.Errorf("<CheckIndex> get indexes: %v", err)
	}
	for _, index := range originIndexs {
		if checkIndexKey(index.Key, key) && unique == index.Unique {
			return nil
		}
	}
	return fmt.Errorf("<CheckIndex> not found index: %v unique: %v", key, unique)
}

func checkIndexKey(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, k := range a {
		if k != b[i] {
			return false
		}
	}
	return true
}

// EnsureSafe `W` and `WMode` only support replset, `WMode`: support version 2.0+
func EnsureSafe(session *mgo.Session, safe *Safe) {
	if safe == nil {
		return
	}
	session.EnsureSafe(&mgo.Safe{
		W:        safe.W,
		WMode:    safe.WMode,
		WTimeout: safe.WTimeout,
		FSync:    safe.FSync,
		J:        safe.J,
	})
}
