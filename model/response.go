package model

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Payload T      `json:"payload,omitempty"`
}

type ApplicationResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_url"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewResponse[T any](message string, payload T) *Response[T] {
	return &Response[T]{
		Message: message,
		Payload: payload,
	}
}
