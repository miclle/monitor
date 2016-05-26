package mgoutil

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
)

var gModes = map[string]int{
	"eventual":  0,
	"monotonic": 1,
	"mono":      1,
	"strong":    2,
}

// SetMode set mgo session mode
func SetMode(s *mgo.Session, modeFriendly string, refresh bool) {

	mode, ok := gModes[strings.ToLower(modeFriendly)]
	if !ok {
		log.Fatal("invalid mgo mode")
	}
	switch mode {
	case 0:
		s.SetMode(mgo.Eventual, refresh)
	case 1:
		s.SetMode(mgo.Monotonic, refresh)
	case 2:
		s.SetMode(mgo.Strong, refresh)
	}
}
