package util

import (
	"encoding/base64"
	"net/http"
)

// BasicTransport a http round tripper
type BasicTransport struct {
	Username  string
	Password  string
	transport http.RoundTripper
}

// RoundTrip : The Request's URL and Header fields must be initialized.
func (t *BasicTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	auth := "Basic " + base64.URLEncoding.EncodeToString([]byte(t.Username+":"+t.Password))
	req.Header.Set("Authorization", auth)
	return t.transport.RoundTrip(req)
}

// NewBasicTransport return a BasicTransport
func NewBasicTransport(username, password string) *BasicTransport {
	return &BasicTransport{
		Username:  username,
		Password:  password,
		transport: http.DefaultTransport,
	}
}
