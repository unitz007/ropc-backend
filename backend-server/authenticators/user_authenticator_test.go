package authenticators

import (
	"backend-server/mocks"
	"backend-server/model/entities"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserAuthenticationFailure(t *testing.T) {

	var userRepositoryMock *mocks.UserRepository

	t.Run("User does not exist", func(t *testing.T) {

		userRepositoryMock = new(mocks.UserRepository)
		userAuthenticator := &userAuthenticator{
			userRepository: userRepositoryMock,
		}

		// should fail if user does not exist
		userRepositoryMock.On("GetUser", mock.Anything).Return(nil, errors.New("invalid user"))
		user, err := userAuthenticator.Authenticate(WrongUsername, mock.Anything)
		assert.Equal(t, err.Error(), "invalid user")
		assert.Nil(t, user)
		userRepositoryMock.AssertCalled(t, "GetUser", mock.Anything)

	})

	t.Run("Wrong password", func(t *testing.T) {

		userRepositoryMock = new(mocks.UserRepository)
		userAuthenticator := &userAuthenticator{
			userRepository: userRepositoryMock,
		}

		// should fail with wrong password combination
		userRepositoryMock.On("GetUser", RightUsername).Return(&entities.User{Username: RightUsername, Password: hashedRightPassword}, nil)
		user, err := userAuthenticator.Authenticate(RightUsername, WrongPassword)
		assert.EqualError(t, err, "invalid user credentials")
		assert.Nil(t, user)
		userRepositoryMock.AssertCalled(t, "GetUser", RightUsername)

	})

	t.Run("Encrypted but wrong password", func(t *testing.T) {

		userRepositoryMock = new(mocks.UserRepository)
		userAuthenticator := &userAuthenticator{
			userRepository: userRepositoryMock,
		}

		// test with right encryption but still wrong password nonetheless, should fail
		userRepositoryMock.On("GetUser", RightUsername).Return(&entities.User{Username: RightUsername, Password: hashedRightPassword}, nil)
		user, err := userAuthenticator.Authenticate(RightUsername, WrongPassword)
		assert.EqualError(t, err, "invalid user credentials")
		assert.Nil(t, user)
		userRepositoryMock.AssertCalled(t, "GetUser", RightUsername)

	})

}

func Test_UserAuthenticationSuccess(t *testing.T) {

	// successful authentication
	userRepositoryMock := new(mocks.UserRepository)
	userAuthenticator := &userAuthenticator{
		userRepository: userRepositoryMock,
	}
	userRepositoryMock.On("GetUser", RightUsername).Return(&entities.User{Username: RightUsername, Password: hashedRightPassword}, nil)
	user, err := userAuthenticator.Authenticate(RightUsername, RightPassword)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	userRepositoryMock.AssertCalled(t, "GetUser", RightUsername)
}
