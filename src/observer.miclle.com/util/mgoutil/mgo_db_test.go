package mgoutil

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
)

func TestIsSessionClosed(t *testing.T) {
	s := Dail("localhost", "strong", 1)
	assert.False(t, IsSessionClosed(s))
	assert.NoError(t, checkSession(s))
	s.Close()
	assert.True(t, IsSessionClosed(s))
	assert.Panics(t, func() { checkSession(s) })
}

func TestCopyDatabase(t *testing.T) {
	db1 := Dail("localhost", "strong", 1).DB("observer_test")
	db2 := CopyDatabase(db1)
	assert.NotEqual(t, db1, db2)
	assert.Equal(t, db1.Name, db2.Name)
	assert.NotEqual(t, db1.Session, db2.Session)
}

func TestCloseDatabase(t *testing.T) {
	db1 := Dail("localhost", "strong", 1).DB("observer_test")
	assert.False(t, IsSessionClosed(db1.Session))
	assert.NoError(t, checkSession(db1.Session))
	CloseDatbase(db1)
	assert.True(t, IsSessionClosed(db1.Session))
	assert.Panics(t, func() { checkSession(db1.Session) })
}

func TestCopyCollection(t *testing.T) {
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	c2 := CopyCollection(c1)
	assert.NotEqual(t, c1, c2)
	assert.NotEqual(t, c1.Database, c2.Database)
	assert.NotEqual(t, c1.Database.Session, c2.Database.Session)
	assert.Equal(t, c1.FullName, c2.FullName)
}

func TestCopyCollectionWhenShutdown(t *testing.T) {
	if os.Getenv("MANUAL_TEST") == "" {
		return
	}
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	for i := 0; i < 10; i++ {
		CopyCollection(c1)
	}
	log.Println("SLEEP STARTED, PLEASE SHUTDOWN MONGODB")
	time.Sleep(time.Second * 10)
	log.Println("SLEEP ENDED")
	c2 := CopyCollection(c1)
	assert.NotEqual(t, c1, c2)
	assert.NotEqual(t, c1.Database, c2.Database)
	assert.NotEqual(t, c1.Database.Session, c2.Database.Session)
	assert.Equal(t, c1.FullName, c2.FullName)
}

func TestCloseCollection(t *testing.T) {
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	assert.False(t, IsSessionClosed(c1.Database.Session))
	assert.NoError(t, checkSession(c1.Database.Session))
	CloseCollection(c1)
	assert.True(t, IsSessionClosed(c1.Database.Session))
	assert.Panics(t, func() { checkSession(c1.Database.Session) })
}

func TestCheckIndex(t *testing.T) {
	c := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	c.DropCollection()

	indexKey := []string{"foo", "bar"}
	noIndexKey := []string{"xxxy"}
	index := mgo.Index{Key: indexKey, Unique: true}

	assert.Nil(t, c.EnsureIndex(index))

	assert.Nil(t, CheckIndex(c, indexKey, true))
	assert.NotNil(t, CheckIndex(c, noIndexKey, true))
	assert.NotNil(t, CheckIndex(c, noIndexKey, false))
	assert.NotNil(t, CheckIndex(c, indexKey, false))
}

/*
Benchmark result:
   BenchmarkCopyCollection10	   20000	     82717 ns/op
   BenchmarkCopyCollection100	   20000	     87755 ns/op
   BenchmarkFastCopyCollection10	 1000000	      1039 ns/op
   BenchmarkFastCopyCollection100	 1000000	      1040 ns/op
*/
func BenchmarkCopyCollection10(b *testing.B) {
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	for i := 0; i < b.N; i += 10 {
		c2 := make([]*mgo.Collection, 10)
		for j := 0; j < 10; j++ {
			c2[j] = CopyCollection(c1)
		}
		for j := 0; j < 10; j++ {
			CloseCollection(c2[j])
		}
	}
}
func BenchmarkCopyCollection100(b *testing.B) {
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	for i := 0; i < b.N; i += 100 {
		c2 := make([]*mgo.Collection, 100)
		for j := 0; j < 100; j++ {
			c2[j] = CopyCollection(c1)
		}
		for j := 0; j < 100; j++ {
			CloseCollection(c2[j])
		}
	}
}
func BenchmarkFastCopyCollection10(b *testing.B) {
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	for i := 0; i < b.N; i += 10 {
		c2 := make([]*mgo.Collection, 10)
		for j := 0; j < 10; j++ {
			c2[j] = FastCopyCollection(c1)
		}
		for j := 0; j < 10; j++ {
			CloseCollection(c2[j])
		}
	}
}
func BenchmarkFastCopyCollection100(b *testing.B) {
	c1 := Dail("localhost", "strong", 1).DB("observer_test").C("coll_test")
	for i := 0; i < b.N; i += 100 {
		c2 := make([]*mgo.Collection, 100)
		for j := 0; j < 100; j++ {
			c2[j] = FastCopyCollection(c1)
		}
		for j := 0; j < 100; j++ {
			CloseCollection(c2[j])
		}
	}
}
