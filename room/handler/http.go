package handler

import (
	"encoding/json"
	"integration-go/api"
	"integration-go/qismo"
	"integration-go/room"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	svc *room.Service
}

func NewHttpHandler(svc *room.Service) *httpHandler {
	return &httpHandler{
		svc: svc,
	}
}

func (h *httpHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		api.WriteResponseJSON(w, http.StatusBadRequest, api.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	room, err := h.svc.GetRoomByID(ctx, int64(id))
	if err != nil {
		api.WriteResponseJSONFromError(w, err)
		return
	}

	api.WriteResponseJSON(w, http.StatusOK, room)
}

func (h *httpHandler) WebhookQismoNewSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req qismo.WebhookNewSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.WriteResponseJSON(w, http.StatusBadRequest, api.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	if err := h.svc.CreateRoom(ctx, &req); err != nil {
		api.WriteResponseJSONFromError(w, err)
		return
	}

	api.WriteResponseJSON(w, http.StatusOK, "ok")
}
