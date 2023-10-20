package services

import (
	"backend-server/mocks"
	model "backend-server/model/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {

	var testUser *model.User
	var userMockRepository *mocks.UserRepository

	t.Run("user without username", func(t *testing.T) {

		userMockRepository = new(mocks.UserRepository)
		testUser = &model.User{
			Username: "",
		}

		userMockRepository.On("CreateUser", testUser).Return(mock.Anything)
		_, err := NewUserService(userMockRepository).CreateUser(testUser)
		assert.EqualError(t, err, "username is required")
		userMockRepository.AssertNotCalled(t, "CreateUser", testUser)
	})

	t.Run("user without password", func(t *testing.T) {
		userMockRepository = new(mocks.UserRepository)
		testUser = &model.User{
			Username: "User",
			Password: "",
		}

		userMockRepository.On("CreateUser", testUser).Return(mock.Anything)
		_, err := NewUserService(userMockRepository).CreateUser(testUser)
		assert.EqualError(t, err, "password is required")
		userMockRepository.AssertNotCalled(t, "CreateUser", testUser)
	})

	t.Run("user without email", func(t *testing.T) {
		userMockRepository = new(mocks.UserRepository)
		testUser = &model.User{
			Username: "User",
			Password: "Password",
			Email:    "",
		}

		userMockRepository.On("CreateUser", testUser).Return(mock.Anything)
		_, err := NewUserService(userMockRepository).CreateUser(testUser)
		assert.EqualError(t, err, "email is required")
		userMockRepository.AssertNotCalled(t, "CreateUser", testUser)

	})

	t.Run("user with already existing username", func(t *testing.T) {

		testUsername := "username"

		testUser = &model.User{
			Username: testUsername,
			Password: "Password",
			Email:    "anpch@example.com",
		}

		testUser2 := &model.User{
			Username: testUsername,
			Password: "Password2",
			Email:    "anpch@example.com",
		}

		userMockRepository = new(mocks.UserRepository)
		userMockRepository.On("GetUserByUsernameOrEmail", testUsername, mock.Anything).Return(testUser, nil)
		_, err := NewUserService(userMockRepository).CreateUser(testUser2)
		assert.EqualError(t, err, "A user with this username already exists")
		userMockRepository.AssertCalled(t, "GetUserByUsernameOrEmail", testUsername, mock.Anything)
	})

	t.Run("user with already existing email", func(t *testing.T) {

		testEmail := "anpch@example.com"

		testUser = &model.User{
			Username: "username2",
			Password: "Password2",
			Email:    testEmail,
		}

		testUserTwo := &model.User{
			Username: "username",
			Password: "Password",
			Email:    testEmail,
		}

		userMockRepository = new(mocks.UserRepository)
		userMockRepository.On("GetUserByUsernameOrEmail", "mock.Anything", testEmail).Return(testUser, nil)
		_, err := NewUserService(userMockRepository).CreateUser(testUserTwo)
		assert.EqualError(t, err, "A user with this email already exists")
		userMockRepository.AssertCalled(t, "GetUserByUsernameOrEmail", "mock.Anything", testEmail)
	})

	t.Run("successful registration", func(t *testing.T) {
		userMockRepository = new(mocks.UserRepository)

		password := "Password"

		testUser = &model.User{
			Username: "User",
			Password: password,
			Email:    "test@gmail.com",
		}

		userMockRepository.On("CreateUser", testUser).Return(testUser, nil)
		userMockRepository.On("GetUserByUsernameOrEmail", testUser.Username, testUser.Email).Return(nil, nil)
		newUser, err := NewUserService(userMockRepository).CreateUser(testUser)
		assert.NoError(t, err, "error is supposed to be nil")
		assert.NotNil(t, newUser, "newUser should be not be nil")
		userMockRepository.AssertCalled(t, "GetUserByUsernameOrEmail", testUser.Username, testUser.Email)
		userMockRepository.AssertCalled(t, "CreateUser", testUser)

		encErr := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(password))
		assert.NoError(t, encErr, "user password was not encrypted")
	})
}
