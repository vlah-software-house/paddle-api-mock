package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

var requestCounter uint64

func NextRequestID() string {
	n := atomic.AddUint64(&requestCounter, 1)
	return fmt.Sprintf("req_%08d", n)
}

type contextKey string

const requestIDKey contextKey = "request_id"

// GetRequestID extracts the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return NextRequestID()
}

// JSON sets Content-Type to application/json on all responses.
func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Auth returns middleware that validates the Bearer token on /v1/ paths.
func Auth(apiKey string, noAuth bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if noAuth || !strings.HasPrefix(r.URL.Path, "/v1/") {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				rid := GetRequestID(r.Context())
				http.Error(w, `{"error":{"type":"authentication_error","code":"unauthorized","detail":"Missing Authorization header"},"meta":{"request_id":"`+rid+`"}}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token == auth || token != apiKey {
				rid := GetRequestID(r.Context())
				http.Error(w, `{"error":{"type":"authentication_error","code":"unauthorized","detail":"Invalid API key"},"meta":{"request_id":"`+rid+`"}}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID injects a request ID into the context and response header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := NextRequestID()
		w.Header().Set("X-Request-ID", rid)
		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
