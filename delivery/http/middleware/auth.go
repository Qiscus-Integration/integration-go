package middleware

import (
	"net/http"
	"os"
)

// ApiKey is a middleware function that checks if the incoming request contains a valid API key
// specified in the Authorization header.
func ApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if apiKey != os.Getenv("API_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
