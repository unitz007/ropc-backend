package kernel

import (
	"fmt"
	"net/http"
	"ropc-backend/utils"
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
	router      Router
	handlers    []api
	middlewares []func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)
}

func (s *server) AttachMiddleware(middleware func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)) Server {
	s.middlewares = append(s.middlewares, middleware)
	return s
}

func NewServer(router Router, defaultMiddlewares []func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)) Server {

	return &server{
		router:      router,
		handlers:    make([]api, 0),
		middlewares: defaultMiddlewares,
	}
}

func (s *server) Start(addr string) error {

	var l Logger = utils.NewZapLogger(utils.NewConfig())

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
		l.Info(fmt.Sprintf("%d handler(s) registered", len(s.handlers)))
		msg := fmt.Sprintf("Server started on port %s, with %s.", PORT, s.router.Name())
		l.Info(msg)
	}()

	err := s.router.Serve(addr)

	if err != nil {
		return err
	}

	return nil

}

func (s *server) RegisterHandler(path, method string, handler func(w http.ResponseWriter, r *http.Request)) {

	for _, m := range s.middlewares {
		m(handler)
	}

	var l Logger = utils.NewZapLogger(utils.NewConfig())

	//newRelicApp := utils.NewRelicInstance().App
	//fHandler := middlewares.RequestLogger(middlewares.PanicRecovery(handler))
	//
	//// register new relic monitor
	//newrelic.WrapHandleFunc(newRelicApp, path, fHandler)
	//txn := s.newRelicApp.StartTransaction(path + "_monitor")
	//defer txn.End()

	switch method {
	case http.MethodGet:
		s.router.Get(path, handler)
	case http.MethodPost:
		s.router.Post(path, handler)
	case http.MethodPut:
		s.router.Put(path, handler)
	case http.MethodDelete:
		s.router.Delete(path, handler)
	default:
		m := fmt.Sprintf("%s not registered: %s", path, fmt.Sprintf("%s is not a upported HTTP method type.", method))
		l.Warn(m)
	}

	h := api{
		path:   path,
		method: method,
	}

	s.handlers = append(s.handlers, h)
}
