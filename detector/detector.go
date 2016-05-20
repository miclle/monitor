package detector

import (
	"net/http"

	"github.com/miclle/observer/xlog"
)

// ProbeGet probe HTTP GET
func ProbeGet(d Delivery, url string) {
	client := &Client{
		Client:   &http.Client{},
		Delivery: d,
	}
	client.Get(xlog.NewLogger(), url)
}
