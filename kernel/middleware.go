package kernel

import (
	"fmt"
	"net/http"
	"ropc-backend/model"
	"ropc-backend/utils"
)

type Middleware interface {
	PanicHandler(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
	RequestLogging(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
}
type middleware struct {
	logger Logger
	config utils.Config
}

func NewMiddleware(logger Logger) Middleware {
	return &middleware{
		logger: logger,
	}
}

func (m middleware) PanicHandler(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorMsg := "Something went wrong"
				if e, ok := err.(error); ok {
					errorMsg = e.Error()
				}
				m.logger.Error(errorMsg)
				_ = utils.PrintResponse(http.StatusBadRequest, w, model.NewResponse[any](errorMsg, nil))
			}
		}()

		h(w, r)
	}
}

func (m middleware) RequestLogging(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		s := fmt.Sprintf("%v request to %v", r.Method, r.URL.Path)
		m.logger.Info(s)
		h(w, r)

	}

}
