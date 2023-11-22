package authenticators

//
//import (

//)
//
//const InvalidClientMessage = "invalid credentials"
//const ConnectionErrorMessage = "could not authenticate client"
//
//type Authenticator interface {
//	Authenticate(clientId, clientSecret string) (*model.AccessToken, error)
//}
//
//type clientAuthenticator struct {
//	repository                    repositories.ApplicationRepository
//	thirdPartyClientAuthenticator ThirdPartyClientAuthenticator
//}
//
//func (selfC clientAuthenticator) Authenticate(clientId, clientSecret string) (*model.AccessToken, error) {
//
//	client, err := selfC.repository.GetByClientId(clientId)
//
//	if err != nil {
//		return nil, err
//	}
//
//	if err = bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(clientSecret)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
//		return nil, errors.New(InvalidClientMessage)
//	}
//
//
//	tokenResponse := &model.TokenResponse{
//		AccessToken: "",
//	}
//
//	return tokenResponse, nil
//}
