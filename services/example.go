package services

import (
	"net/http"

	"github.com/miclle/observer/detector"
	"github.com/miclle/observer/xlog"
)

// ProbeGithubAPI probe github API
// Github API: (GET https://api.github.com)
func ProbeGithubAPI(delivery detector.Delivery) {
	client := &detector.Client{
		Client:   &http.Client{},
		Delivery: delivery,
	}
	client.Get(xlog.NewLogger(), "https://api.github.com")
}
