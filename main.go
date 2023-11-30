//   Product Api:
//    version: 0.1
//    title: Product Api
//   Schemes: http, https
//   Host:
//   BasePath: /api/v1
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   SecurityDefinitions:
//    Bearer:
//     type: apiKey
//     name: Authorization
//     in: header
//   swagger:meta
package main

import (
	"log"
	"net/http"
	"ropc-backend/handlers"
	"ropc-backend/kernel"
	"ropc-backend/middlewares"
	"ropc-backend/repositories"
	"ropc-backend/services"
	"ropc-backend/utils"

	"gorm.io/gorm"
)

const (
	loginPath          = "/token"
	appPath            = "/apps"
	generateSecretPath = "/apps/generate_secret"
	userPath           = "/users"
)

func main() {

	config := utils.NewConfig()
	ctx, err := kernel.NewContext[gorm.DB](config)
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	applicationRepository := repositories.NewApplicationRepository(ctx.Database)
	userRepository := repositories.NewUserRepository(ctx.Database)

	// services
	authenticatorService := services.NewAuthenticatorService(applicationRepository, config)

	// Handlers
	authenticationHandler := handlers.NewAuthenticationHandler(authenticatorService)
	applicationHandler := handlers.NewApplicationHandler(applicationRepository, ctx.Router)
	userHandler := handlers.NewUserHandler(config, userRepository)

	security := middlewares.NewSecurity(config, userRepository)

	// Server
	server := kernel.NewServer(ctx.Router)
	server.RegisterHandler(appPath, http.MethodPost, security.TokenValidation(applicationHandler.CreateApplication))
	server.RegisterHandler(appPath, http.MethodGet, security.TokenValidation(applicationHandler.GetApplications))
	server.RegisterHandler(appPath+"/{client_id}", http.MethodGet, security.TokenValidation(applicationHandler.GetApplication))
	server.RegisterHandler(appPath+"/{client_id}", http.MethodDelete, security.TokenValidation(applicationHandler.DeleteApplication))

	server.RegisterHandler(generateSecretPath, http.MethodPut, security.TokenValidation(applicationHandler.GenerateSecret))
	server.RegisterHandler(loginPath, http.MethodPost, authenticationHandler.Authenticate)
	server.RegisterHandler(userPath, http.MethodPost, userHandler.CreateUser)
	server.RegisterHandler(userPath+"/auth", http.MethodPost, userHandler.AuthenticateUser)

	// swagger

	ctx.Logger.Fatal(server.Start(":" + config.ServerPort()).Error())
}
