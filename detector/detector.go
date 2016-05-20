package detector

import (
	"net/http"

	"github.com/miclle/observer/util"
	"github.com/miclle/observer/xlog"
)

// ProbeGet probe HTTP GET request
func ProbeGet(d Delivery, url string) {
	client := &Client{
		Client:   &http.Client{},
		Delivery: d,
	}
	client.Get(xlog.NewLogger(), url)
}

// ProbeBasicAuthGet probe HTTP Basic Authentication GET request
func ProbeBasicAuthGet(d Delivery, url, username, password string) {
	tr := util.NewBasicTransport(username, password)
	client := &Client{
		Client:   &http.Client{Transport: tr},
		Delivery: d,
	}
	client.Get(xlog.NewLogger(), url)
}
