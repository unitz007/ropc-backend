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

func NewServer(ctx Context, defaultMiddlewares Middleware) Server {

	return &server{
		Context:     ctx,
		handlers:    make([]api, 0),
		middlewares: defaultMiddlewares,
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

	fHandler := s.middlewares.RequestLogging(s.middlewares.PanicHandler(handler))

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
