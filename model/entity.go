package model

import (
	"errors"

	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	ClientId     string `gorm:"index;unique;not-null"`
	ClientSecret string `gorm:"size:100"`
	Name         string `gorm:"unique"`
	RedirectUri  string `gorm:""`
}

func NewApplication(clientId, name string) (*Application, error) {

	if clientId == "" {
		return nil, errors.New("client id should not be empty")
	}

	app := &Application{
		ClientId: clientId,
		Name:     name,
	}

	return app, nil
}

func (a *Application) ToDTO() *ApplicationDto {
	return &ApplicationDto{
		ClientId: a.ClientId,
		Name:     a.Name,
	}
}
