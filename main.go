package main

import (
	"backend-server/handlers"
	"backend-server/middlewares"
	"backend-server/repositories"
	"backend-server/routers"
	"backend-server/services"
	"backend-server/utils"
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
	userPath           = "/users"
)

func main() {

	config := utils.NewConfig()
	var router routers.Router

	switch config.Mux() {
	case "gorilla_mux":
		router = routers.NewRouter(mux.NewRouter())
	case "chi_router":
		router = routers.NewChiRouter(chi.NewRouter())
	default:
		s := fmt.Sprintf("%s is not a supported Mux router", config.Mux())
		utils.NewLogger().Error(s, true)
	}

	DB := repositories.NewDataBase(config)

	// Repositories
	applicationRepository := repositories.NewApplicationRepository(DB)
	userRepository := repositories.NewUserRepository(DB)

	// services
	authenticatorService := services.NewAuthenticatorService(applicationRepository, config)

	// Handlers
	authenticationHandler := handlers.NewAuthenticationHandler(authenticatorService)
	applicationHandler := handlers.NewApplicationHandler(applicationRepository, router)
	userHandler := handlers.NewUserHandler(config, userRepository)

	security := middlewares.NewSecurity(config, userRepository)

	// Server
	server := NewServer(router)
	server.RegisterHandler(appPath, http.MethodPost, security.TokenValidation(applicationHandler.CreateApplication))
	server.RegisterHandler(appPath, http.MethodGet, security.TokenValidation(applicationHandler.GetApplications))
	server.RegisterHandler(appPath+"/{client_id}", http.MethodGet, security.TokenValidation(applicationHandler.GetApplication))
	server.RegisterHandler(generateSecretPath, http.MethodPut, security.TokenValidation(applicationHandler.GenerateSecret))
	server.RegisterHandler(loginPath, http.MethodPost, authenticationHandler.Authenticate)
	server.RegisterHandler(userPath, http.MethodPost, userHandler.CreateUser)
	server.RegisterHandler(userPath+"/auth", http.MethodPost, userHandler.AuthenticateUser)

	log.Fatal(server.Start(":" + config.ServerPort()))
}
