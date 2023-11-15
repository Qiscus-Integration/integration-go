package httpclient

import (
	"net/http"
	"time"
)

type retryTransport struct {
	base    http.RoundTripper
	retries int
}

func (rt *retryTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var attempts int

	for attempts < rt.retries {
		resp, err = rt.base.RoundTrip(req)
		attempts++

		if err == nil && resp.StatusCode != 500 && resp.StatusCode != 429 {
			return
		}

		time.Sleep(1 << rt.retries)
	}

	return

}
