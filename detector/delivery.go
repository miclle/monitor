package detector

import (
	"net/http"
	"time"
)

// Delivery func after the completion of the request triggering method
type Delivery func(l Logger, start, end time.Time, req *http.Request, resp *http.Response, err error)
