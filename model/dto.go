package model

type ApplicationDto struct {
	ClientId    string `json:"client_id"`
	Name        string `json:"name"`
	RedirectURL string `json:"redirect_url"`
}
