package handler

import (
	"encoding/json"
	"integration-go/httpresp"
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
		httpresp.WriteFail(w, http.StatusBadRequest, httpresp.Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	room, err := h.svc.GetRoomByID(ctx, int64(id))
	if err != nil {
		httpresp.WriteFailFromError(w, err)
		return
	}

	httpresp.WriteSuccess(w, http.StatusOK, room)

}

func (h *httpHandler) WebhookQismoNewSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req qismo.WebhookNewSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httpresp.WriteFail(w, http.StatusBadRequest, httpresp.Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	if err := h.svc.CreateRoom(ctx, &req); err != nil {
		httpresp.WriteFailFromError(w, err)
		return
	}

	httpresp.WriteSuccess(w, http.StatusOK, "ok")
}
