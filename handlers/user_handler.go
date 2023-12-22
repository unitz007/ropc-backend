package handlers

import (
	"errors"
	"net/http"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	userCreated = "User created Successfully"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	AuthenticateUser(w http.ResponseWriter, r *http.Request)
	GetUserDetails(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	config         utils.Config
	userRepository kernel.Repository[model.User]
}

func (u *userHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest *model.LoginRequest

	err := JsonToStruct(r.Body, &loginRequest)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	condition := utils.Queries[utils.WhereUsernameOrEmailIs](loginRequest.UsernameOrEmail)

	user, err := u.userRepository.Get(condition)
	if err != nil {
		_ = utils.PrintResponse[any](http.StatusUnauthorized, w, nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		_ = utils.PrintResponse[any](http.StatusUnauthorized, w, nil)
		return
	}

	accessToken := model.NewAccessToken(user.Email, u.config.TokenExpiry(), utils.GetIssuerUri(u.config)).Sign(u.config.TokenSecret())

	resp := model.NewResponse[*model.TokenResponse]("Authentication successful", &model.TokenResponse{AccessToken: accessToken})

	_ = utils.PrintResponse[*model.Response[*model.TokenResponse]](http.StatusOK, w, resp)
}

func NewUserHandler(config utils.Config, userRepository kernel.Repository[model.User]) UserHandler {
	return &userHandler{
		config:         config,
		userRepository: userRepository,
	}
}

func (u *userHandler) CreateUser(response http.ResponseWriter, request *http.Request) {

	var requestBody *model.CreateUser

	err := JsonToStruct(request.Body, &requestBody)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	if requestBody.UserName == utils.Blank {
		panic(errors.New("username is required"))
	}

	if requestBody.EmailAddress == utils.Blank {
		panic(errors.New("email is required"))
	}

	if requestBody.Password == utils.Blank {
		panic(errors.New("password is required"))
	}

	user := model.User{
		Username: requestBody.UserName,
		Password: requestBody.Password,
		Email:    requestBody.EmailAddress,
	}

	err = u.userRepository.Create(user)
	if err != nil {
		message := err.Error()
		switch {
		case strings.Contains(err.Error(), "username"):
			message = "username already exists"
		case strings.Contains(err.Error(), "email"):
			message = "email already exists"
		}

		panic(kernel.NewError(http.StatusConflict, message))
	}

	res := model.NewResponse[any](userCreated, nil)

	_ = utils.PrintResponse(http.StatusCreated, response, res)
}

func (u *userHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		accessToken := r.Header.Get("Authorization")
		claims, err := utils.ValidateToken(accessToken, utils.NewConfig().TokenSecret())
		if err != nil {
			panic(err)
		}

		userDetails := model.UserDetails{
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			ClientId: claims["client_id"].(string),
		}

		_ = utils.PrintResponse(http.StatusOK, w, model.NewResponse("User details fetched successfully",
			userDetails))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
