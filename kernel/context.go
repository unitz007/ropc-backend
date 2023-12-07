package kernel

import (
	"errors"
	"ropc-backend/routers"
	"ropc-backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

type Context interface {
	Database() Database
	Router() Router
	Logger() Logger
}

type context struct {
	db     Database
	router Router
	logger Logger
}

func (c context) Database() Database {
	return c.db
}

func (c context) Router() Router {
	return c.router
}

func (c context) Logger() Logger {
	return c.logger
}

func NewContext(config utils.Config) (Context, error) {

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
			return nil, errors.New("invalid router specified: " + config.Mux())
		}
	}()

	if err != nil {
		return nil, err
	}

	context := context{
		db:     db,
		logger: utils.NewZapLogger(config),
		router: router,
	}

	//fmt.Println("Application context loaded successfully")
	//fmt.Printf("_______________________________________\n")
	//fmt.Printf("Logger ==> Zap Logger                 |\n")
	//fmt.Printf("Router ==> %s\n", context.router.Name())
	//fmt.Printf("Database ==> %s\n", context.Database()..Name())
	//fmt.Printf("---------------------------------------\n")

	return context, nil
}
