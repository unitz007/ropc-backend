package kernel

import (
	"errors"
	"ropc-backend/routers"
	"ropc-backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

type Context[DatabaseConnectionReference any] struct {
	Database Database[DatabaseConnectionReference]
	Router   Router
	Logger   Logger
}

func NewContext[DatabaseConnectionReference any](config utils.Config) (*Context[DatabaseConnectionReference], error) {

	db, err := NewDatabase(config)
	if err != nil {
		return nil, errors.New("database connection error: " + err.Error())
	}

	router, err := func() (Router, error) {
		switch config.Mux() {
		case "chi_router":
			return routers.NewChiRouter(chi.NewRouter()), nil
		case "gorilla_mux":
			return routers.NewRouter(mux.NewRouter()), nil
		default:
			return nil, errors.New("invalid router specified:" + config.Mux())
		}
	}()

	context := &Context[DatabaseConnectionReference]{
		Database: db,
		Logger:   utils.NewZapLogger(config),
		Router:   router,
	}

	return context, nil
}
