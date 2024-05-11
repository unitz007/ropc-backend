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
	"ropc-backend/model"
	"ropc-backend/services"
	"ropc-backend/utils"

	"gorm.io/gorm"
)

const (
	loginPath          = "/token"
	appPath            = "/apps"
	generateSecretPath = "/apps/generate_secret"
	userPath           = "/users"
	idkPath            = "/idk"
)

func main() {

	config := utils.NewConfig()
	ctx, err := kernel.NewContext(config)
	if err != nil {
		log.Fatal(err)
	}

	defaultMiddlewares := kernel.NewMiddleware(ctx)

	db, ok := ctx.Database().GetDatabaseConnection().(*gorm.DB)
	if !ok {
		log.Fatal("Database connection failed. could not get Database connection object")
	}

	// Repositories
	applicationRepository := kernel.NewRepository[model.Application](model.Application{}, db)
	userRepository := kernel.NewRepository[model.User](model.User{}, db)

	// services
	authenticatorService := services.NewAuthenticatorService(applicationRepository, config)

	// Handlers
	authenticationHandler := handlers.NewAuthenticationHandler(authenticatorService, ctx)
	applicationHandler := handlers.NewApplicationHandler(applicationRepository, ctx)
	userHandler := handlers.NewUserHandler(config, userRepository)
	authorizationHandler := handlers.NewAuthorizationHandler()

	security := kernel.NewSecurity(config, userRepository)

	// Server
	server := kernel.NewServer(ctx, defaultMiddlewares, security)

	server.RegisterHandler(appPath, http.MethodPost, applicationHandler.CreateApplication, true)
	server.RegisterHandler(appPath, http.MethodGet, applicationHandler.GetApplications, true)
	server.RegisterHandler(appPath+"/{client_id}", http.MethodGet, applicationHandler.GetApplication, true)
	server.RegisterHandler(appPath+"/{client_id}", http.MethodDelete, applicationHandler.DeleteApplication, true)
	server.RegisterHandler(generateSecretPath, http.MethodPut, applicationHandler.GenerateSecret, true)
	server.RegisterHandler(loginPath, http.MethodPost, authenticationHandler.Authenticate, false)
	server.RegisterHandler(idkPath, http.MethodGet, authorizationHandler.IDK, false)

	server.RegisterHandler(userPath, http.MethodPost, userHandler.CreateUser, false)

	// swagger

	ctx.Logger().Fatal(server.Start(":" + config.ServerPort()).Error())
}
