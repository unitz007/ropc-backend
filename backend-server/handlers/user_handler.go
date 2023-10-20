package handlers

import (
	"backend-server/conf"
	"backend-server/model"
	"backend-server/repositories"
	"backend-server/services"
	"backend-server/utils"
	"net/http"
)

const (
	userCreated = "User created Successfully"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserDetails(w http.ResponseWriter, r *http.Request)
}

type userHandler struct{}

func NewUserAuthenticatorHandler() UserHandler {
	return &userHandler{}
}

func (u *userHandler) CreateUser(response http.ResponseWriter, request *http.Request) {
	var user *model.User

	err := JsonToStruct(request.Body, &user)
	if err != nil {
		panic(err)
	}

	database := conf.NewDataBase(conf.EnvironmentConfig)

	userRepository := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepository)

	_, err = userService.CreateUser(user)
	if err != nil {
		panic(err)

	}
	_ = PrintResponse(http.StatusCreated, response, model.NewResponse[any](userCreated, nil))
}

func (u *userHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		accessToken := r.Header.Get("Authorization")
		claims, err := utils.ValidateToken(accessToken, conf.EnvironmentConfig.TokenSecret())
		if err != nil {
			panic(err)
		}

		userDetails := model.UserDetails{
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			ClientId: claims["client_id"].(string),
		}

		_ = PrintResponse(http.StatusOK, w, model.NewResponse("User details fetched successfully",
			userDetails))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
