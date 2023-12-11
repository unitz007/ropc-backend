package services

import (
	"errors"
	"fmt"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

const (
	InvalidClientMessage = "invalid credentials"
)

type AuthenticatorService interface {
	ClientCredentials(clientId, clientSecret string) (string, error)
}

type authenticatorService struct {
	kernel.Context
	applicationRepository kernel.Repository[model.Application]
	config                utils.Config
}

func (a *authenticatorService) ClientCredentials(clientId, clientSecret string) (string, error) {

	condition := utils.Queries[utils.WhereClientIdIs](clientId)

	app, err := a.applicationRepository.Get(condition)
	if err != nil {
		return utils.Blank, err
	}

	fmt.Println(app)

	if err = bcrypt.CompareHashAndPassword([]byte(app.ClientSecret), []byte(clientSecret)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return utils.Blank, errors.New(InvalidClientMessage)
	}

	accessToken := model.NewAccessToken(app.ClientId, a.config.TokenExpiry(), utils.GetIssuerUri(a.config)).Sign(a.config.TokenSecret())
	if err != nil {
		return utils.Blank, err
	}

	return accessToken, nil
}

func NewAuthenticatorService(applicationRepository kernel.Repository[model.Application], config utils.Config) AuthenticatorService {
	return &authenticatorService{applicationRepository: applicationRepository, config: config}
}
