package authenticators

type ThirdPartyClientAuthenticator interface {
	Authenticate(clientID string, clientSecret string) (*bool, error)
}
