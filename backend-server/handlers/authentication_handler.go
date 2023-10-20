package handlers

import (
	"backend-server/model"
	"backend-server/services"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

const authenticationSuccessMsg = "Authentication successful"

type AuthenticationHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type authenticationHandler struct {
	authenticator services.AuthenticatorService
}

func NewAuthenticationHandler(authenticator services.AuthenticatorService) AuthenticationHandler {
	return &authenticationHandler{authenticator}
}

func (a *authenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		panic(errors.New("invalid content-type"))
	}

	clientId := r.FormValue("client_id")
	if clientId == "" {
		panic(errors.New("client id is required"))
	}

	clientSecret := r.FormValue("client_secret")
	if clientSecret == "" {
		panic(errors.New("client secret is required"))
	}

	grantType := r.FormValue("grant_type")
	if grantType == "" {
		panic(errors.New("grant type is required"))
	}

	var token *model.Token
	var err error

	switch grantType {
	case "client_credentials":
		token, err = a.authenticator.ClientCredentials(clientId, clientSecret)
	default:
		panic(errors.New("invalid grant type"))
	}

	if err != nil {
		_ = PrintResponse(http.StatusUnauthorized, w, &model.Response[string]{Message: err.Error()})
		return
	}

	response := &model.Response[*model.Token]{
		Message: authenticationSuccessMsg,
		Payload: token,
	}

	_ = PrintResponse[*model.Response[*model.Token]](http.StatusOK, w, response)
}

func (a *authenticationHandler) GetMux() *mux.Router {
	return &mux.Router{}
}
