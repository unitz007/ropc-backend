package services

import (
	"errors"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/repositories"
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
	applicationRepository repositories.ApplicationRepository
	config                utils.Config
}

func (a *authenticatorService) ClientCredentials(clientId, clientSecret string) (string, error) {
	app, err := a.applicationRepository.GetByClientId(clientId)
	if err != nil {
		return utils.Blank, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(app.ClientSecret), []byte(clientSecret)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return utils.Blank, errors.New(InvalidClientMessage)
	}

	accessToken := model.NewAccessToken(app.ClientId, a.config.TokenExpiry(), utils.GetIssuerUri(a.config)).Sign(a.config.TokenSecret())
	if err != nil {
		return utils.Blank, err
	}

	return accessToken, nil
}

func NewAuthenticatorService(applicationRepository repositories.ApplicationRepository, config utils.Config) AuthenticatorService {
	return &authenticatorService{applicationRepository: applicationRepository, config: config}
}
