package resp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type customError struct {
	code    int
	message string
}

func (e customError) Error() string {
	return e.message
}

func (e customError) HTTPStatusCode() int {
	return e.code
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		data         interface{}
		expectedBody string
		expectedCode int
	}{
		{
			name:         "Write simple data",
			code:         http.StatusOK,
			data:         map[string]string{"message": "success"},
			expectedBody: `{"message":"success"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Write empty struct",
			code:         http.StatusCreated,
			data:         Empty{},
			expectedBody: `{}`,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, tt.code, tt.data)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestWriteJSONWithPaginate(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		data         interface{}
		total        int
		page         int
		limit        int
		expectedCode int
		expectedMeta Meta
	}{
		{
			name:         "Single page of data",
			code:         http.StatusOK,
			data:         []string{"item1", "item2"},
			total:        2,
			page:         1,
			limit:        10,
			expectedCode: http.StatusOK,
			expectedMeta: Meta{
				Page:      1,
				PageTotal: 1,
				Total:     2,
			},
		},
		{
			name:         "Multiple pages",
			code:         http.StatusOK,
			data:         []string{"item1", "item2"},
			total:        15,
			page:         2,
			limit:        5,
			expectedCode: http.StatusOK,
			expectedMeta: Meta{
				Page:      2,
				PageTotal: 3,
				Total:     15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONWithPaginate(w, tt.code, tt.data, tt.total, tt.page, tt.limit)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedCode, w.Code)

			var response DataPaginate
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMeta, response.Meta)
		})
	}
}

func TestWriteJSONFromError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Standard error",
			err:          errors.New("standard error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error_code":500,"error_message":"Something went wrong"}`,
		},
		{
			name: "Custom HTTP error",
			err: customError{
				code:    http.StatusBadRequest,
				message: "invalid input",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error_code":400,"error_message":"invalid input"}`,
		},
		{
			name:         "Nil error",
			err:          nil,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error_code":500,"error_message":"Something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONFromError(w, tt.err)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestMeta_JSON(t *testing.T) {
	meta := Meta{
		Page:      2,
		PageTotal: 5,
		Total:     50,
	}

	jsonData, err := json.Marshal(meta)
	assert.NoError(t, err)

	expected := `{"page":2,"page_total":5,"total":50}`
	assert.JSONEq(t, expected, string(jsonData))
}

func TestDataPaginate_JSON(t *testing.T) {
	data := DataPaginate{
		Data: []string{"item1", "item2"},
		Meta: Meta{
			Page:      1,
			PageTotal: 1,
			Total:     2,
		},
	}

	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	expected := `{
		"data": ["item1", "item2"],
		"meta": {
			"page": 1,
			"page_total": 1,
			"total": 2
		}
	}`
	assert.JSONEq(t, expected, string(jsonData))
}

func TestHTTPError_JSON(t *testing.T) {
	httpError := HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid input",
	}

	jsonData, err := json.Marshal(httpError)
	assert.NoError(t, err)

	expected := `{"error_code":400,"error_message":"invalid input"}`
	assert.JSONEq(t, expected, string(jsonData))
}
