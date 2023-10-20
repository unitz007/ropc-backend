package main

import (
	"backend-server/logger"
	"backend-server/middlewares"
	"backend-server/routers"
	"fmt"
	"net/http"
	"time"
)

type Server interface {
	Start(address string) error
	RegisterHandler(path, method string, handler func(w http.ResponseWriter, r *http.Request))
	AttachMiddleware(middleware func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)) Server
}

type api struct {
	path   string
	method string
}

type server struct {
	router      routers.Router
	handlers    []api
	middlewares []func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)
}

func (s *server) AttachMiddleware(middleware func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)) Server {
	s.middlewares = append(s.middlewares, middleware)
	return s
}

func NewServer(router routers.Router) Server {
	return &server{
		//config:   config,
		router:   router,
		handlers: make([]api, 0),
	}
}

func (s *server) Start(addr string) error {

	PORT := func() string {
		index := 0
		for i, v := range addr {
			if v == ':' {
				index = i
			}
		}

		return addr[index+1:]

	}()

	go func() {
		time.Sleep(time.Millisecond * 5)
		logger.Info(fmt.Sprintf("%d handler(s) registered", len(s.handlers)))
		msg := fmt.Sprintf("Server started on port %s, with %s.", PORT, s.router.Name())
		logger.Info(msg)
	}()

	err := s.router.Serve(addr)

	if err != nil {
		return err
	}

	return nil

}

func (s *server) RegisterHandler(path, method string, handler func(w http.ResponseWriter, r *http.Request)) {

	fHandler := middlewares.RequestLogger(middlewares.PanicRecovery(handler))

	switch method {
	case http.MethodGet:
		s.router.Get(path, fHandler)
	case http.MethodPost:
		s.router.Post(path, fHandler)
	case http.MethodPut:
		s.router.Put(path, fHandler)
	default:
		m := fmt.Sprintf("%s not registered: %s", path, fmt.Sprintf("%s is an unsupported method type.", method))
		logger.Warn(m)
	}

	h := api{
		path:   path,
		method: method,
	}

	s.handlers = append(s.handlers, h)
}
