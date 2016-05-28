package mgoutil

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoIndexService_EnsureIndex(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if !assert.NoError(t, err) {
		return
	}
	defer session.Close()

	collection := session.DB("mgo_util_test").C("test")
	_ = collection.DropCollection()
	service := NewMongoIndexService(session, "mgo_util_test")
	service.EnsureIndex(bson.M{
		"name": "test",
		"unique": []string{
			"uid",
			"email",
		},
	})
	indexes, err := collection.Indexes()
	assert.NoError(t, err)
	assert.Equal(t, len(indexes), 3)
}
