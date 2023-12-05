package kernel

import (
	"fmt"
	"net/http"
	"time"
)

type Server interface {
	Start(address string) error
	RegisterHandler(path, method string, handler func(w http.ResponseWriter, r *http.Request))
	//AttachMiddleware(middleware func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)) Server
}

type api struct {
	path   string
	method string
}

type server struct {
	Context
	handlers    []api
	middlewares Middleware
}

//func (s *server) AttachMiddleware(middleware func(w http.ResponseWriter, r *http.Request) func(http.ResponseWriter, *http.Request)) Server {
//	//s.middlewares = append(s.middlewares, middleware)
//	return s
//}

func NewServer(ctx Context, defaultMiddlewares Middleware) Server {

	return &server{
		Context:     ctx,
		handlers:    make([]api, 0),
		middlewares: defaultMiddlewares,
	}
}

func (s *server) Start(addr string) error {

	//var l Logger = utils.NewZapLogger(utils.NewConfig())

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
		s.Logger().Info(fmt.Sprintf("%d handler(s) registered", len(s.handlers)))
		msg := fmt.Sprintf("Server started on port %s, with %s.", PORT, s.Router().Name())
		s.Logger().Info(msg)
	}()

	err := s.Router().Serve(addr)

	if err != nil {
		return err
	}

	return nil

}

func (s *server) RegisterHandler(path, method string, handler func(w http.ResponseWriter, r *http.Request)) {

	//for _, m := range s.middlewares {
	//	h := m(handler)
	//}

	//var l Logger = utils.NewZapLogger(utils.NewConfig())

	//newRelicApp := utils.NewRelicInstance().App
	fHandler := s.middlewares.RequestLogging(s.middlewares.PanicHandler(handler))
	//
	//// register new relic monitor
	//newrelic.WrapHandleFunc(newRelicApp, path, fHandler)
	//txn := s.newRelicApp.StartTransaction(path + "_monitor")
	//defer txn.End()

	switch method {
	case http.MethodGet:
		s.Router().Get(path, fHandler)
	case http.MethodPost:
		s.Router().Post(path, fHandler)
	case http.MethodPut:
		s.Router().Put(path, fHandler)
	case http.MethodDelete:
		s.Router().Delete(path, fHandler)
	default:
		m := fmt.Sprintf("%s not registered: %s", path, fmt.Sprintf("%s is not a upported HTTP method type.", method))
		s.Logger().Warn(m)
	}

	h := api{
		path:   path,
		method: method,
	}

	s.handlers = append(s.handlers, h)
}
