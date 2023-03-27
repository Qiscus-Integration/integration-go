package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type apiQismo struct {
	appID     string
	secretKey string
}

// NewApiQismo creates and returns a new instance of the `apiQismo` struct which implements the `OmnichannelRepository` interface.
// It takes two string arguments: `appID` and `secretKey`, which are used to authenticate with the Qismo API.
func NewApiQismo(appID, secretKey string) *apiQismo {
	return &apiQismo{appID, secretKey}
}

func (q *apiQismo) CreateRoomTag(ctx context.Context, roomID string, tag string) (err error) {
	url := "https://multichannel.qiscus.com/api/v1/room_tag/create"
	payload, _ := json.Marshal(map[string]interface{}{
		"room_id": roomID,
		"tag":     tag,
	})

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Qiscus-App-Id", q.appID)
	req.Header.Add("Qiscus-Secret-Key", q.secretKey)

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("api %s returned status code %d. response body:%s", res.Request.URL.String(), res.StatusCode, string(resBody))
		return
	}

	return
}

func (q *apiQismo) ResolvedRoom(ctx context.Context, roomID string) (err error) {
	url := "https://multichannel.qiscus.com/api/v1/admin/service/mark_as_resolved"
	payload, _ := json.Marshal(map[string]interface{}{
		"room_id": roomID,
	})

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Qiscus-App-Id", q.appID)
	req.Header.Add("Qiscus-Secret-Key", q.secretKey)

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("api %s returned status code %d. response body:%s", res.Request.URL.String(), res.StatusCode, string(resBody))
		return
	}

	return
}
