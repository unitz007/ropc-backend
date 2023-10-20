package services

import (
	"backend-server/model"
	"backend-server/repositories"
	"backend-server/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	InvalidClientMessage   = "invalid credentials"
	ConnectionErrorMessage = "could not authenticate client"
)

type AuthenticatorService interface {
	ClientCredentials(clientId, clientSecret string) (*model.Token, error)
}

type authenticatorService struct {
	applicationRepository repositories.ApplicationRepository
}

func (a *authenticatorService) ClientCredentials(clientId, clientSecret string) (*model.Token, error) {
	app, err := a.applicationRepository.GetByClientId(clientId)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(app.ClientSecret), []byte(clientSecret)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, errors.New(InvalidClientMessage)
	}

	tokenResponse, err := utils.GenerateToken(app, clientSecret)
	if err != nil {
		return nil, err
	}

	return &model.Token{AccessToken: tokenResponse}, nil
}

func NewAuthenticatorService(applicationRepository repositories.ApplicationRepository) AuthenticatorService {
	return &authenticatorService{applicationRepository: applicationRepository}
}

//
//func NewAuthenticator(cA authenticators.ClientAuthenticator) Authenticator {
//	return &authenticator{
//		//userAuthenticator:   uA,
//		clientAuthenticator: cA,
//	}
//}
//
//func (selfC authenticator) Authenticate(client *entities.Application) (*dto.Token, error) {
//
//	c, err := selfC.clientAuthenticator.Authenticate(client.ClientId, client.ClientSecret)
//	if err != nil {
//		return nil, err
//	}
//
//	accessToken, err := utils.GenerateToken(c, conf.EnvironmentConfig.TokenSecret())
//	if err != nil {
//		return nil, err
//	}
//
//	token := &dto.Token{AccessToken: accessToken}
//
//	return token, nil
//}
