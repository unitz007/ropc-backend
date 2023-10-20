package authenticators

import (
	"backend-server/model"
	"backend-server/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticator interface {
	Authenticate(username, password string) (*model.User, error)
}

type userAuthenticator struct {
	userRepository repositories.UserRepository
}

func NewUserAuthenticator(userRepository repositories.UserRepository) UserAuthenticator {
	return &userAuthenticator{
		userRepository: userRepository,
	}
}

func (selfC userAuthenticator) Authenticate(username, password string) (*model.User, error) {

	user, err := selfC.userRepository.GetUser(username)

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, errors.New("invalid user credentials")
	}

	return user, nil
}
