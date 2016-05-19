package services

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/miclle/observer/detector"
	"github.com/stretchr/testify/assert"
)

var delivery = func(l detector.Logger, start, end time.Time, req *http.Request, resp *http.Response, err error) {

	latency := end.Sub(start)

	var code int

	if resp != nil {
		code = resp.StatusCode
	}

	log.Printf("Probe Github API: %s, StatusCode: %d, latency: %13v", req.URL.String(), code, latency)
}

func TestProbeGithubAPI(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(delivery)
	ProbeGithubAPI(delivery)
}
