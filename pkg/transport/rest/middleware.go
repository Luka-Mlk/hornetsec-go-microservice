package rest

import (
	"context"
	"document-metadata/pkg/constants"
	"net/http"

	"github.com/google/uuid"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(constants.HeaderRequestID)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Set(constants.HeaderRequestID, id)

		ctx := context.WithValue(r.Context(), constants.HeaderRequestID, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
