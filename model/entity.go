package model

import (
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	ClientId     string `gorm:"index;unique;not-null"`
	ClientSecret string `gorm:"size:100"`
	Name         string `gorm:""`
	RedirectUri  string `gorm:""`
	UserID       uint
	User         User `gorm:"references:ID"`
}

type User struct {
	gorm.Model
	Username string `gorm:"index;unique;not-null"`
	Password string `gorm:"not-null"`
	Email    string `gorm:"not-null;unique"`
}

type Test struct {
	gorm.Model
	TestValue string `gorm:""`
}

func (a Application) ToDTO() *ApplicationDto {
	return &ApplicationDto{
		ClientId:    a.ClientId,
		Name:        a.Name,
		RedirectURL: a.RedirectUri,
	}
}
