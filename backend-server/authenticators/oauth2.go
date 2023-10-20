package authenticators

import "backend-server/model"

type Oauth2 interface {
	ClientCredentials(clientId, clientSecret string) (*model.Token, error)
}

type GrantType interface {
}
