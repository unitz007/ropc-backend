package routers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type MuxMultiplexer struct {
	router *mux.Router
}

func (m *MuxMultiplexer) GetPathVariable(req *http.Request, variable string) (error, string) {
	params := mux.Vars(req)
	p := params[variable]
	if p == "" {
		return errors.New("missing path variable: " + variable), ""
	}

	return nil, p
}

func (m *MuxMultiplexer) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.router.HandleFunc(path, handler).Methods(http.MethodDelete)
}

func (m *MuxMultiplexer) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.router.HandleFunc(path, handler).Methods(http.MethodPut)
}

func (m *MuxMultiplexer) Name() string {
	return "Mux Router"
}

func NewRouter(mux *mux.Router) Router {
	return &MuxMultiplexer{router: mux}
}

func (m *MuxMultiplexer) Get(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	m.router.HandleFunc(path, handlerFunc).Methods(http.MethodGet)
}

func (m *MuxMultiplexer) Serve(addr string) error {
	return http.ListenAndServe(addr, m.router)
}

func (m *MuxMultiplexer) Post(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	m.router.HandleFunc(path, handlerFunc).Methods(http.MethodPost)
}
