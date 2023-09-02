package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type retryRoundTripper struct {
	roundTripper http.RoundTripper
	ctx          context.Context
	retries      int
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var attempts int

	for attempts < rt.retries {
		resp, err := rt.roundTripper.RoundTrip(req)
		attempts++

		if err == nil && resp.StatusCode != 500 && resp.StatusCode != 429 {
			return resp, err
		}

		log.Ctx(rt.ctx).Info().Msgf(
			"%s %s request failed (retry %d/%d)",
			req.Method,
			req.URL,
			attempts,
			rt.retries,
		)

		time.Sleep(1 << rt.retries)
	}

	return nil, fmt.Errorf("request failed after all retries")

}

func MakeHTTPRequest(ctx context.Context, method, url string, body []byte, headers map[string]string, response interface{}) (err error) {
	client := &http.Client{
		Timeout: time.Second * 20,
		Transport: &retryRoundTripper{
			roundTripper: http.DefaultTransport,
			ctx:          ctx,
			retries:      3,
		},
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"%s %s http client returned error %d response: %s",
			resp.Request.Method,
			resp.Request.URL,
			resp.StatusCode,
			string(responseBody))
	}

	if response != nil {
		if err = json.Unmarshal(responseBody, response); err != nil {
			return
		}
	}

	return
}
