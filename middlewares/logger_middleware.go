package middlewares

import (
	"fmt"
	"net/http"
	"ropc-backend/utils"
)

func RequestLogger(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		s := fmt.Sprintf("%v request to %v", r.Method, r.URL.Path)
		utils.NewLogger().Info(s)
		h(w, r)

	}

}
