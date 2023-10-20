package main

import (
	"backend-server/conf"
	"backend-server/handlers"
	"backend-server/logger"
	"backend-server/repositories"
	"backend-server/routers"
	"backend-server/services"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

const (
	loginPath          = "/token"
	appPath            = "/apps"
	generateSecretPath = "/apps/generate_secret"
)

func main() {

	config := conf.EnvironmentConfig
	var router routers.Router

	switch config.Mux() {
	case "gorilla_mux":
		router = routers.NewRouter(mux.NewRouter())
	case "chi_router":
		router = routers.NewChiRouter(chi.NewRouter())
	default:
		s := fmt.Sprintf("%s is not a supported Mux router", config.Mux())
		logger.Error(s, true)
	}

	DB := conf.NewDataBase(config)

	// Repositories
	applicationRepository := repositories.NewApplicationRepository(DB)

	// services
	authenticatorService := services.NewAuthenticatorService(applicationRepository)

	// Handlers
	authenticationHandler := handlers.NewAuthenticationHandler(authenticatorService)
	applicationHandler := handlers.NewApplicationHandler(applicationRepository)

	// Server
	server := NewServer(router)
	server.RegisterHandler(appPath, http.MethodPost, applicationHandler.CreateApplication)
	server.RegisterHandler(appPath, http.MethodGet, applicationHandler.GetApplications)
	server.RegisterHandler(generateSecretPath, http.MethodPut, applicationHandler.GenerateSecret)
	server.RegisterHandler(loginPath, http.MethodPost, authenticationHandler.Authenticate)
	server.RegisterHandler(appPath, http.MethodGet, applicationHandler.GetApplication)

	log.Fatal(server.Start(":" + config.ServerPort()))
}
