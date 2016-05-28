package detector_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/miclle/observer/detector"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Reply
func Reply(w http.ResponseWriter, code int, data interface{}) {

	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	h := w.Header()
	h.Set("Content-Length", strconv.Itoa(len(msg)))
	h.Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msg)
}

var userAgentTest string

func foo(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	Reply(w, 200, map[string]interface{}{
		"info":  "Call method foo",
		"url":   req.RequestURI,
		"query": req.Form,
	})
}

func agent(w http.ResponseWriter, req *http.Request) {
	userAgentTest = req.Header.Get("User-Agent")
}

type Object struct {
}

func (p *Object) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqBytes, _ := ioutil.ReadAll(req.Body)
	Reply(w, 200, map[string]interface{}{"info": "Call method object", "req": string(reqBytes)})
}

var done = make(chan bool)

func server(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", foo)
	mux.Handle("object", new(Object))
	return httptest.NewServer(mux)
}

func TestDo(t *testing.T) {

	s := httptest.NewServer(http.HandlerFunc(agent))
	defer s.Close()

	c := detector.DefaultClient
	{
		req, _ := http.NewRequest("GET", s.URL+"/agent", nil)
		c.Do(nil, req)
		assert.Equal(t, userAgentTest, "Observer detector package")
	}
	{
		req, _ := http.NewRequest("GET", s.URL+"/agent", nil)
		req.Header.Set("User-Agent", "tst")
		c.Do(nil, req)
		assert.Equal(t, userAgentTest, "tst")
	}
}
