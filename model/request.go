package model

type CreateApplication struct {
	ClientId    string `json:"client_id"`
	Name        string `json:"name"`
	RedirectUri string `json:"redirect_uri"`
}
