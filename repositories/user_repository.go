package repositories

import (
	"backend-server/model"
	"errors"
	"log"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(usernameOrEmail string) (*model.User, error)
	GetUserByUsernameOrEmail(username, email string) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
}

type userRepository struct {
	db Database[gorm.DB]
}

func NewUserRepository(database Database[gorm.DB]) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (selfC userRepository) GetUser(username string) (*model.User, error) {
	var user *model.User

	err := selfC.db.GetDatabaseConnection().Model(&model.User{}).Where("username = ? OR email = ?", username, username).First(&user).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("invalid user credentials")
	}

	return user, nil
}

func (selfC userRepository) CreateUser(user *model.User) (*model.User, error) {

	err := selfC.db.GetDatabaseConnection().Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "username") {
			return nil, errors.New("username already exists")
		} else if strings.Contains(err.Error(), "email") {
			return nil, errors.New("email already exists")
		}

		return nil, err
	}

	return user, nil
}

func (selfC userRepository) GetUserByUsernameOrEmail(username, email string) (*model.User, error) {
	var user model.User

	err := selfC.db.GetDatabaseConnection().Where("username = ? OR email = ?", username, email).First(&user).Error

	if err != nil {
		return nil, errors.New("could not execute query")
	}

	return &user, nil
}
