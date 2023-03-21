package http

import (
	"encoding/json"
	"integration-go/delivery/http/middleware"
	"integration-go/domain"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type room struct {
	roomUC domain.RoomUsecase
}

// NewRoom creates a new instance of the room struct with the provided room usecase.
func NewRoom(roomUC domain.RoomUsecase) *room {
	return &room{
		roomUC: roomUC,
	}
}

type newSessionWebhookReq struct {
	IsNewSession bool `json:"is_new_session"`
	Payload      struct {
		Room struct {
			ID              string `json:"id"`
			IDStr           string `json:"id_str"`
			IsPublicChannel bool   `json:"is_public_channel"`
			Name            string `json:"name"`
			Options         string `json:"options"`
			Participants    []struct {
				Email string `json:"email"`
			} `json:"participants"`
			RoomAvatar string `json:"room_avatar"`
			TopicID    string `json:"topic_id"`
			TopicIDStr string `json:"topic_id_str"`
			Type       string `json:"type"`
		} `json:"room"`
	} `json:"payload"`
	WebhookType string `json:"webhook_type"`
}

func (h *room) HandleNewSessionWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body newSessionWebhookReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		room := &domain.Room{
			MultichannelRoomID: body.Payload.Room.IDStr,
		}

		if err := h.roomUC.CreateRoom(ctx, room); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("ok"))

	}
}

func (h *room) FetchRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rooms, err := h.roomUC.FetchRoom(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, _ := json.Marshal(rooms)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	}
}

func (h *room) HandleRoute(r chi.Router) {
	r.Post("/wh/multichannel/new_session", h.HandleNewSessionWebhook())

	r.Group(func(r chi.Router) {
		r.Use(middleware.ApiKey)
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/rooms", h.FetchRoom())
		})
	})
}
