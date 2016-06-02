package detector

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"qiniupkg.com/x/xlog.v7"
)

// UserAgent Info
var UserAgent = "Observer detector package"

// --------------------------------------------------------------------

// Delivery func after the completion of the request triggering method
type Delivery func(l *xlog.Logger, start, end time.Time, req *http.Request, resp *http.Response, err error)

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
func (r Client) deliver(l *xlog.Logger, start time.Time, req *http.Request, resp *http.Response, err error) {
	end := time.Now()
	if r.Delivery != nil {
		r.Delivery(l, start, end, req, resp, err)
	}
}

// --------------------------------------------------------------------

// Head send http head request
func (r Client) Head(l *xlog.Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Get send http get request
func (r Client) Get(l *xlog.Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Delete send http delete request
func (r Client) Delete(l *xlog.Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Post send http post request, no Content-Type set
func (r Client) Post(l *xlog.Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// PostWith send http post request
func (r Client) PostWith(l *xlog.Logger, url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	req.ContentLength = int64(bodyLength)
	return r.Do(l, req)
}

// PostWithForm send http post form request
func (r Client) PostWithForm(l *xlog.Logger, url1 string, data map[string][]string) (resp *http.Response, err error) {

	msg := url.Values(data).Encode()
	return r.PostWith(l, url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

// Do the http request
func (r Client) Do(l *xlog.Logger, req *http.Request) (resp *http.Response, err error) {
	start := time.Now()

	if l != nil {
		req.Header.Set("X-Reqid", l.ReqId)
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
