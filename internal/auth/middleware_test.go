package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMiddleware(t *testing.T) {
	secretKey := "test-secret"
	m := NewMiddleware(secretKey)

	assert.Equal(t, secretKey, m.secretKey, "NewMiddleware() returned wrong secretKey")
}

func TestStaticToken(t *testing.T) {
	tests := []struct {
		name           string
		secretKey      string
		providedToken  string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			secretKey:      "secret123",
			providedToken:  "secret123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid token",
			secretKey:      "secret123",
			providedToken:  "wrongtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty token",
			secretKey:      "secret123",
			providedToken:  "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiddleware(tt.secretKey)
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.providedToken)

			rr := httptest.NewRecorder()
			handler := m.StaticToken(nextHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")

			if tt.expectedStatus == http.StatusUnauthorized {
				expected := `{"error_code":401,"error_message":"Unauthorized"}`
				assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
			}
		})
	}
}
