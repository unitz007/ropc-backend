package model

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Payload T      `json:"payload,omitempty"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

func (t AccessToken) Valid() error {
	return nil
}

type AccessToken struct {
	Issuer    string `json:"iss"`
	Sub       string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Name      string `json:"name"`
}

type ApplicationResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewResponse[T any](message string, payload T) *Response[T] {
	return &Response[T]{
		Message: message,
		Payload: payload,
	}
}
