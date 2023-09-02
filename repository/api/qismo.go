package api

import (
	"context"
	"encoding/json"
	"integration-go/util"
	"net/http"
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

func (q *apiQismo) headers() map[string]string {
	return map[string]string{
		"Qiscus-App-Id":     q.appID,
		"Qiscus-Secret-Key": q.secretKey,
	}
}

func (q *apiQismo) CreateRoomTag(ctx context.Context, roomID string, tag string) (err error) {
	url := "https://multichannel.qiscus.com/api/v1/room_tag/create"
	payload, _ := json.Marshal(map[string]interface{}{
		"room_id": roomID,
		"tag":     tag,
	})

	err = util.MakeHTTPRequest(ctx, http.MethodPost, url, payload, q.headers(), nil)
	return
}

func (q *apiQismo) ResolvedRoom(ctx context.Context, roomID string) (err error) {
	url := "https://multichannel.qiscus.com/api/v1/admin/service/mark_as_resolved"
	payload, _ := json.Marshal(map[string]interface{}{
		"room_id": roomID,
	})

	err = util.MakeHTTPRequest(ctx, http.MethodPost, url, payload, q.headers(), nil)
	return
}
