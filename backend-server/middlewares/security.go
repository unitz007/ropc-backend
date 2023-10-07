package middlewares

import (
	"net/http"
	"ropc-service/handlers"
	"ropc-service/model"
)

const (
	tokenHeader         = "Authorization"
	tokenHeaderErrorMsg = "Bearer token is required"
)

func Security(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(tokenHeader)
		if accessToken == "" {
			_ = handlers.PrintResponse(http.StatusForbidden, w, model.NewResponse[any](tokenHeaderErrorMsg, nil))
		}

		h(w, r)
	}
}
