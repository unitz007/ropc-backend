package services

import (
	"backend-server/model"
	"backend-server/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *model.User) (*model.User, error)
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	return &userService{repository: repository}
}

func (selfC userService) CreateUser(user *model.User) (*model.User, error) {
	if user.Username == "" {
		return nil, errors.New("username is required")
	}

	if user.Password == "" {
		return nil, errors.New("password is required")
	}

	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	existingUser, _ := selfC.repository.GetUserByUsernameOrEmail(user.Username, user.Email)

	if existingUser != nil {
		if existingUser.Username == user.Username {
			return nil, errors.New("A user with this username already exists")
		}

		if existingUser.Email == user.Email {
			return nil, errors.New("A user with this email already exists")
		}

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	user.Password = string(hashedPassword)

	newUser, err := selfC.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
