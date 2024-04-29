package server

import (
	"integration-go/api/resp"
	"net/http"
)

// ApiKey is a middleware function that checks if the incoming request contains a valid API key
// specified in the Authorization header.
func staticTokenAuthMiddleware(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tokenStr := r.Header.Get("Authorization")
			if tokenStr != secretKey {
				resp.WriteJSON(w, http.StatusUnauthorized, resp.HTTPError{
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
				})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
