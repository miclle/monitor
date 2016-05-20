package detector

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var delivery = func(l Logger, start, end time.Time, req *http.Request, resp *http.Response, err error) {

	latency := end.Sub(start)

	var code int

	if resp != nil {
		code = resp.StatusCode
	}

	log.Printf("Probe HTTP GET:%s, StatusCode: %d, latency: %13v", req.URL.String(), code, latency)
}

func TestProbeBasicAuthGet(t *testing.T) {
	assert := assert.New(t)
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	assert.NotEmpty(username, "BASIC_AUTH username must require")
	assert.NotEmpty(password, "BASIC_AUTH password must require")
}
