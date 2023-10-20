package middlewares

import (
	"backend-server/logger"
	"fmt"
	"net/http"
)

func RequestLogger(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		s := fmt.Sprintf("%v request to %v", r.Method, r.URL.Path)
		logger.Info(s)
		h(w, r)

	}

}
