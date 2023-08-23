package routes

import (
	"net/http"

	"github.com/NumexaHQ/monger/model"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r = r.WithContext(model.AssignRequestID(ctx))

		h.ServeHTTP(w, r)
	})
}
