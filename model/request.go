package model

// swagger:model CreateApplication
type CreateApplication struct {

	// ClientId of the application
	// in: string
	ClientId string `json:"client_id"`

	// Name of application
	// in: string
	Name string `json:"name"`

	// RedirectUri of application
	// in: string
	RedirectUri string `json:"redirect_uri"`
}

// swagger:model CreateUser
type CreateUser struct {

	// Username of User
	// in: string
	UserName string `json:"username"`
	// EmailAddress of User
	// in: string
	EmailAddress string `json:"email"`
	// Password of User
	// in: string
	Password string `json:"password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}
