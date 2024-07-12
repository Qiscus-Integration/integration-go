package room

import (
	"encoding/json"
	"integration-go/internal/api/resp"
	"integration-go/internal/qismo"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type httpHandler struct {
	svc *Service
}

func NewHttpHandler(svc *Service) *httpHandler {
	return &httpHandler{
		svc: svc,
	}
}

func (h *httpHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, resp.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	room, err := h.svc.GetRoomByID(ctx, int64(id))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed to get room: %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, room)
}

func (h *httpHandler) WebhookQismoNewSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req qismo.WebhookNewSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, resp.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	if err := h.svc.CreateRoom(ctx, &req); err != nil {
		log.Ctx(ctx).Error().Msgf("failed to create room: %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, "ok")
}
