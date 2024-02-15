package api

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
)

type Meta struct {
	Page      int `json:"page"`
	PageTotal int `json:"page_total"`
	Total     int `json:"total"`
}

type DataPaginate struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type HTTPError struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error"`
}

type Empty struct{}

func WriteResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func WriteResponseJSONWithPaginate(w http.ResponseWriter, statusCode int, data interface{}, total int, page int, limit int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := Meta{
		Page:      page,
		PageTotal: totalPages,
		Total:     total,
	}

	jsonData, _ := json.Marshal(DataPaginate{
		Data: data,
		Meta: meta,
	})

	w.Write(jsonData)
}

func WriteResponseJSONFromError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := "Something went wrong"

	var httpErr interface{ HTTPStatusCode() int }
	if errors.As(err, &httpErr) {
		code = httpErr.HTTPStatusCode()
		msg = err.Error()
	}

	errResp := HTTPError{
		StatusCode: code,
		Message:    msg,
	}

	resp, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
