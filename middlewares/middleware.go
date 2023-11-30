package middlewares

import (
	"net/http"
	"ropc-backend/handlers"
	"ropc-backend/kernel"
	"ropc-backend/model"
)

type Middleware struct {
	logger kernel.Logger
}

func NewMiddleware(logger kernel.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m Middleware) PanicRecovery(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorMsg := "Something went wrong"
				if e, ok := err.(error); ok {
					errorMsg = e.Error()
				}
				m.logger.Error(errorMsg)
				_ = handlers.PrintResponse(http.StatusBadRequest, w, model.NewResponse[any](errorMsg, nil))
			}
		}()

		h(w, r)
	}
}
