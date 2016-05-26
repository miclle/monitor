package detector

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// UserAgent Info
var UserAgent = "Observer detector package"

// --------------------------------------------------------------------

// Client is HTTP Client
type Client struct {
	*http.Client
	Delivery Delivery
}

// DefaultClient is a default HTTP client
var DefaultClient = Client{Client: &http.Client{Transport: DefaultTransport}}

// NewClientTimeout return a timeout HTTP client
func NewClientTimeout(dial, resp time.Duration, delivery Delivery) Client {
	return Client{
		Client:   &http.Client{Transport: NewTransportTimeout(dial, resp)},
		Delivery: delivery,
	}
}

// --------------------------------------------------------------------

// deliver call the delivery method after the http do method execute complete
func (r Client) deliver(l Logger, start time.Time, req *http.Request, resp *http.Response, err error) {
	end := time.Now()
	if r.Delivery != nil {
		r.Delivery(l, start, end, req, resp, err)
	}
}

// --------------------------------------------------------------------

// Logger Interface
type Logger interface {
	XRequestID() string
}

// --------------------------------------------------------------------

// Head send http head request
func (r Client) Head(l Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Get send http get request
func (r Client) Get(l Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Delete send http delete request
func (r Client) Delete(l Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// PostEx send http post request, no Content-Type set
func (r Client) PostEx(l Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// PostWith send http post request
func (r Client) PostWith(l Logger, url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(bodyLength)
	return r.Do(l, req)
}

// PostWith64 send http post request
func (r Client) PostWith64(l Logger, url1 string, bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(l, req)
}

// PostWithForm send http post form request
func (r Client) PostWithForm(l Logger, url1 string, data map[string][]string) (resp *http.Response, err error) {

	msg := url.Values(data).Encode()
	return r.PostWith(l, url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

// PostWithJSON send http post request, Content-Type: application/json
func (r Client) PostWithJSON(l Logger, url1 string, data interface{}) (resp *http.Response, err error) {

	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.PostWith(l, url1, "application/json", bytes.NewReader(msg), len(msg))
}

// Do the http request
func (r Client) Do(l Logger, req *http.Request) (resp *http.Response, err error) {
	start := time.Now()

	if l != nil {
		req.Header.Set("X-Request-Id", l.XRequestID())
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}

	resp, err = r.Client.Do(req)

	if err != nil {
		r.deliver(l, start, req, resp, err)
		return
	}

	r.deliver(l, start, req, resp, err)
	return
}
