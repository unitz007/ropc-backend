package handlers

import (
	"errors"
	"net/http"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/services"
	"ropc-backend/utils"
)

const (
	authenticationSuccessMsg      = "Authentication successful"
	grantTypeParam                = "grant_type"
	clientSecretParam             = "client_secret"
	clientIdParam                 = "client_id"
	contentTypeHeader             = "Content-Type"
	contentTypeErrorMessage       = "invalid content-type"
	clientIdErrorMessage          = "client id is required"
	clientSecretErrorMessage      = "client secret is required"
	grantTypeRequiredErrorMessage = "grant type is required"
	invalidGrantTypeErrorMessage  = "invalid grant type"
)

type AuthenticationHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type authenticationHandler struct {
	kernel.Context
	authenticator services.AuthenticatorService
}

func NewAuthenticationHandler(authenticator services.AuthenticatorService, ctx kernel.Context) AuthenticationHandler {
	return &authenticationHandler{ctx, authenticator}
}

func (a *authenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get(contentTypeHeader) != "application/x-www-form-urlencoded" {
		panic(errors.New(contentTypeErrorMessage))
	}

	clientId := r.FormValue(clientIdParam)
	if clientId == utils.Blank {
		panic(errors.New(clientIdErrorMessage))
	}

	clientSecret := r.FormValue(clientSecretParam)
	if clientSecret == utils.Blank {
		panic(errors.New(clientSecretErrorMessage))
	}

	grantType := r.FormValue(grantTypeParam)
	if grantType == utils.Blank {
		panic(errors.New(grantTypeRequiredErrorMessage))
	}

	var token string
	var err error

	switch grantType {
	case "client_credentials":
		token, err = a.authenticator.ClientCredentials(clientId, clientSecret)
	default:
		panic(errors.New(invalidGrantTypeErrorMessage))
	}

	if err != nil {
		_ = utils.PrintResponse(http.StatusUnauthorized, w, &model.Response[string]{Message: err.Error()})
		return
	}

	tokenResponse := &model.TokenResponse{AccessToken: token}

	response := &model.Response[*model.TokenResponse]{
		Message: authenticationSuccessMsg,
		Payload: tokenResponse,
	}

	_ = utils.PrintResponse[*model.Response[*model.TokenResponse]](http.StatusOK, w, response)
}
