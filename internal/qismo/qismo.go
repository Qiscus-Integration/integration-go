package qismo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Qismo struct {
	client    httpClient
	url       string
	appID     string
	secretKey string
}

func New(client httpClient, url, appID, secretKey string) *Qismo {
	return &Qismo{
		client:    client,
		url:       url,
		appID:     appID,
		secretKey: secretKey,
	}
}

func (q *Qismo) headers() map[string]string {
	return map[string]string{
		"Qiscus-App-Id":     q.appID,
		"Qiscus-Secret-Key": q.secretKey,
	}
}

func (q *Qismo) CreateRoomTag(ctx context.Context, roomID string, tag string) error {
	url := fmt.Sprintf("%s/api/v1/room_tag/create", q.url)
	payload, _ := json.Marshal(map[string]any{
		"room_id": roomID,
		"tag":     tag,
	})

	err := q.client.Call(ctx, http.MethodPost, url, bytes.NewBuffer(payload), q.headers(), nil)
	return err
}

func (q *Qismo) ResolvedRoom(ctx context.Context, roomID string) error {
	url := fmt.Sprintf("%s/api/v1/admin/service/mark_as_resolved", q.url)
	payload, _ := json.Marshal(map[string]any{
		"room_id": roomID,
	})

	err := q.client.Call(ctx, http.MethodPost, url, bytes.NewBuffer(payload), q.headers(), nil)
	return err
}
