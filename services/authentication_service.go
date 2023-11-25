package services

import (
	"errors"
	"ropc-backend/model"
	"ropc-backend/repositories"
	"ropc-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

const (
	InvalidClientMessage   = "invalid credentials"
	ConnectionErrorMessage = "could not authenticate client"
)

type AuthenticatorService interface {
	ClientCredentials(clientId, clientSecret string) (string, error)
}

type authenticatorService struct {
	applicationRepository repositories.ApplicationRepository
	config                utils.Config
}

func (a *authenticatorService) ClientCredentials(clientId, clientSecret string) (string, error) {
	app, err := a.applicationRepository.GetByClientId(clientId)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(app.ClientSecret), []byte(clientSecret)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", errors.New(InvalidClientMessage)
	}

	accessToken := model.NewAccessToken(app.ClientId, a.config).Sign()
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func NewAuthenticatorService(applicationRepository repositories.ApplicationRepository, config utils.Config) AuthenticatorService {
	return &authenticatorService{applicationRepository: applicationRepository, config: config}
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
