package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpClient struct {
	client *http.Client
}

func New(client *http.Client) *HttpClient {
	return &HttpClient{
		client: client,
	}
}

func (h *HttpClient) Call(ctx context.Context, method, url string, body io.Reader, headers map[string]string, response interface{}) (err error) {
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), url, body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
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
