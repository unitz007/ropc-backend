package handlers

import (
	"backend-server/model"
	"backend-server/repositories"
	"backend-server/utils"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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
	userRepository repositories.UserRepository
}

func (u *userHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest *model.LoginRequest

	err := JsonToStruct(r.Body, &loginRequest)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	user, err := u.userRepository.GetUser(loginRequest.UsernameOrEmail)
	if err != nil {
		_ = PrintResponse[any](http.StatusUnauthorized, w, nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		_ = PrintResponse[any](http.StatusUnauthorized, w, nil)
		return
	}

	accessToken := model.AccessToken{
		Issuer:    "",
		Sub:       user.Email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Duration(u.config.TokenExpiry()) * time.Minute).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken).SignedString([]byte(u.config.TokenSecret()))

	resp := model.NewResponse[*model.Token]("Authentication successful", &model.Token{AccessToken: token})

	_ = PrintResponse[*model.Response[*model.Token]](http.StatusOK, w, resp)
}

func NewUserHandler(config utils.Config, userRepository repositories.UserRepository) UserHandler {
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

	if requestBody.UserName == "" {
		panic(errors.New("username is required"))
	}

	if requestBody.EmailAddress == "" {
		panic(errors.New("email is required"))
	}

	if requestBody.Password == "" {
		panic(errors.New("password is required"))
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 0)

	user := &model.User{
		Username: requestBody.UserName,
		Password: string(hashed),
		Email:    requestBody.EmailAddress,
	}

	_, err = u.userRepository.CreateUser(user)
	if err != nil {
		panic(err)

	}

	res := model.NewResponse[any](userCreated, nil)

	_ = PrintResponse(http.StatusCreated, response, res)
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

		_ = PrintResponse(http.StatusOK, w, model.NewResponse("User details fetched successfully",
			userDetails))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
