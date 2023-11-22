package authenticators

import "backend-server/model"

type Oauth2 interface {
	ClientCredentials(clientId, clientSecret string) (*model.AccessToken, error)
}

type GrantType interface {
}
