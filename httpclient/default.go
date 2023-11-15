package httpclient

import (
	"net/http"
	"time"
)

var DefaultClient = &http.Client{
	Timeout: time.Second * 20,
	Transport: &retryTransport{
		base:    http.DefaultTransport,
		retries: 3,
	},
}
