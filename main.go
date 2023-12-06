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
	"context"
	"errors"
	"log"
	"net/http"
	"ropc-backend/handlers"
	"ropc-backend/kernel"
	"ropc-backend/repositories"
	"ropc-backend/services"
	"ropc-backend/utils"
)

const (
	loginPath           = "/token"
	appPath             = "/apps"
	generateSecretPath  = "/apps/generate_secret"
	userPath            = "/users"
	tokenHeader         = "Authorization"
	tokenHeaderErrorMsg = "bearer token is required"
)

func main() {

	config := utils.NewConfig()
	ctx, err := kernel.NewContext(config)
	if err != nil {
		log.Fatal(err)
	}

	defaultMiddlewares := kernel.NewMiddleware(ctx)
	// Repositories
	applicationRepository := repositories.NewApplicationRepository(ctx.Database())
	userRepository := repositories.NewUserRepository(ctx.Database())

	// services
	authenticatorService := services.NewAuthenticatorService(applicationRepository, config)

	// Handlers
	authenticationHandler := handlers.NewAuthenticationHandler(authenticatorService, ctx)
	applicationHandler := handlers.NewApplicationHandler(applicationRepository, ctx)
	userHandler := handlers.NewUserHandler(config, userRepository)

	security := func(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get(tokenHeader)

			if accessToken == utils.Blank {
				panic(errors.New(tokenHeaderErrorMsg + " for path: " + r.URL.String()))
			}

			token, err := utils.ValidateToken(accessToken, config.TokenSecret())

			if err != nil {
				panic(errors.New("token validation failed: " + err.Error()))
			}

			email := token["sub"].(string)
			user, err := userRepository.GetUser(email)
			if err != nil {
				http.Error(w, "", http.StatusForbidden)
			}

			r = r.WithContext(context.WithValue(r.Context(), handlers.UserKey, user))

			h(w, r)
		}
	}

	// Server
	server := kernel.NewServer(ctx, defaultMiddlewares)
	server.RegisterHandler(appPath, http.MethodPost, security(applicationHandler.CreateApplication))
	server.RegisterHandler(appPath, http.MethodGet, security(applicationHandler.GetApplications))
	server.RegisterHandler(appPath+"/{client_id}", http.MethodGet, security(applicationHandler.GetApplication))
	server.RegisterHandler(appPath+"/{client_id}", http.MethodDelete, security(applicationHandler.DeleteApplication))

	server.RegisterHandler(generateSecretPath, http.MethodPut, security(applicationHandler.GenerateSecret))
	server.RegisterHandler(loginPath, http.MethodPost, authenticationHandler.Authenticate)
	server.RegisterHandler(userPath, http.MethodPost, userHandler.CreateUser)
	server.RegisterHandler(userPath+"/auth", http.MethodPost, userHandler.AuthenticateUser)

	// swagger

	ctx.Logger().Fatal(server.Start(":" + config.ServerPort()).Error())
}
