package qismo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"integration-go/httpclient"
	"net/http"
)

type client struct {
	url       string
	appID     string
	secretKey string
}

// NewApiQismo creates and returns a new instance of the `apiQismo` struct which implements the `OmnichannelRepository` interface.
// It takes two string arguments: `appID` and `secretKey`, which are used to authenticate with the Qismo API.
func NewClient(url, appID, secretKey string) *client {
	return &client{
		url:       url,
		appID:     appID,
		secretKey: secretKey,
	}
}

func (q *client) headers() map[string]string {
	return map[string]string{
		"Qiscus-App-Id":     q.appID,
		"Qiscus-Secret-Key": q.secretKey,
	}
}

func (q *client) CreateRoomTag(ctx context.Context, roomID string, tag string) error {
	url := fmt.Sprintf("%s/api/v1/room_tag/create", q.url)
	payload, _ := json.Marshal(map[string]interface{}{
		"room_id": roomID,
		"tag":     tag,
	})

	hc := httpclient.New(http.DefaultClient)
	err := hc.Call(ctx, http.MethodPost, url, bytes.NewBuffer(payload), q.headers(), nil)

	return err
}

func (q *client) ResolvedRoom(ctx context.Context, roomID string) error {
	url := fmt.Sprintf("%s/api/v1/admin/service/mark_as_resolved", q.url)
	payload, _ := json.Marshal(map[string]interface{}{
		"room_id": roomID,
	})

	hc := httpclient.New(http.DefaultClient)
	err := hc.Call(ctx, http.MethodPost, url, bytes.NewBuffer(payload), q.headers(), nil)

	return err
}
