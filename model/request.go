package model

type CreateApplication struct {
	ClientId    string `json:"client_id"`
	Name        string `json:"name"`
	RedirectUri string `json:"redirect_uri"`
}

type CreateUser struct {
	UserName     string `json:"username"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}
