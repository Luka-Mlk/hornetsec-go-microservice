package rest

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", id)

		ctx := context.WithValue(r.Context(), "X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
