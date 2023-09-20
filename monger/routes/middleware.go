package routes

import (
	"net/http"

	"github.com/NumexaHQ/monger/model"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r = r.WithContext(model.AssignRequestID(ctx))

		// validate api key
		apiKey := r.Header.Get("X-Numexa-Api-Key")
		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
