package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"index;unique;not-null"`
	Password string `gorm:"not-null"`
	Email    string `gorm:"not-null;unique"`
}
