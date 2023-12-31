package kernel

import (
	"fmt"
	"net/http"
	"ropc-backend/utils"
	"time"
)

type Middleware interface {
	PanicHandler(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
	RequestLogging(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
}
type middleware struct {
	Context
}

func NewMiddleware(ctx Context) Middleware {
	return &middleware{
		ctx,
	}
}

func (m middleware) PanicHandler(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				code := http.StatusInternalServerError
				errorMsg := "Internal Server Error"
				if e, ok := err.(Error); ok {
					errorMsg = e.Error()
					code = e.Code()
				} else {
					if e, ok := err.(error); ok {
						code = http.StatusBadRequest
						errorMsg = e.Error()
					}
				}

				m.Logger().Error(errorMsg)

				_ = utils.PrintResponseNew[any](w, code, errorMsg, nil)
			}
		}()

		h(w, r)
	}
}

func (m middleware) RequestLogging(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		h(w, r) // executed
		endTime := time.Since(startTime).Milliseconds()
		s := fmt.Sprintf("%v request to %v completed in %dms", r.Method, r.URL.Path, endTime)
		m.Logger().Info(s)

	}

}
