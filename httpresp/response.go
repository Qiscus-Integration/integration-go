package httpresp

import (
	"encoding/json"
	"errors"
	"integration-go/entity"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"error_code"`
	Message    string `json:"error_message"`
}

type Empty struct{}

func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func WriteFail(w http.ResponseWriter, statusCode int, error interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(error)
	w.Write(jsonData)
}

func WriteFailFromError(w http.ResponseWriter, err error) {
	var statusCode int
	if errors.Is(err, entity.ErrNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, entity.ErrInternal) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, entity.ErrDatabase) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, entity.ErrBadRequest) {
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, entity.ErrCantProceed) {
		statusCode = http.StatusUnprocessableEntity
	} else if errors.Is(err, entity.ErrUnauthorized) {
		statusCode = http.StatusUnauthorized
	} else if errors.Is(err, entity.ErrForbidden) {
		statusCode = http.StatusForbidden
	} else {
		statusCode = http.StatusInternalServerError
	}

	errorMsg := Error{
		Message:    err.Error(),
		StatusCode: statusCode,
	}

	resp, _ := json.Marshal(errorMsg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}
