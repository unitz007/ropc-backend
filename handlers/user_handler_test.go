package handlers

import (
	"backend-server/mocks"
	"backend-server/model"
	"backend-server/utils"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	username := "username"
	email := "test@gmail.com"
	password := "password"

	user := &model.User{
		Username: username,
		Password: password,
		Email:    email,
	}

	tt := []struct {
		Test         string
		Request      *strings.Reader
		PanicMessage string
	}{
		{
			Test:         "should panic with 'username is required'",
			Request:      userRequest(t, "", email, password),
			PanicMessage: "username is required",
		},

		{
			Test:         "should panic with 'email is required'",
			Request:      userRequest(t, username, "", password),
			PanicMessage: "email is required",
		},

		{
			Test:         "should panic with 'password is required'",
			Request:      userRequest(t, username, email, ""),
			PanicMessage: "password is required",
		},
	}

	for _, e := range tt {
		userRepository := new(mocks.UserRepository)
		request := httptest.NewRequest(http.MethodPost, "/users", e.Request)
		response := httptest.NewRecorder()

		handler := NewUserHandler(utils.NewConfig(), userRepository)
		exec := func() {
			handler.CreateUser(response, request)
		}

		t.Run(e.Test, func(t *testing.T) {
			assert.PanicsWithError(t, e.PanicMessage, exec)
		})

		userRepository.AssertNotCalled(t, "CreateUser", user)

	}

	t.Run("should not create a new user if username already exists in database", func(t *testing.T) {
		userRepository := new(mocks.UserRepository)
		json := userRequest(t, username, email, password)

		userRepository.On("CreateUser", user).Return(nil, errors.New("user with this username already exists"))

		request := httptest.NewRequest(http.MethodPost, "/users", json)
		response := httptest.NewRecorder()

		handler := NewUserHandler(utils.NewConfig(), userRepository)

		exec := func() {
			handler.CreateUser(response, request)
		}

		assert.PanicsWithError(t, "user with this username already exists", exec)
		userRepository.AssertCalled(t, "CreateUser", user)
	})

	t.Run("user created successfully", func(t *testing.T) {
		userRepository := new(mocks.UserRepository)
		json := userRequest(t, username, email, password)

		userRepository.On("CreateUser", user).Return(user, nil)

		request := httptest.NewRequest(http.MethodPost, "/users", json)
		response := httptest.NewRecorder()

		handler := NewUserHandler(utils.NewConfig(), userRepository)
		handler.CreateUser(response, request)

		expected := http.StatusCreated
		got := response.Code

		if expected != got {
			t.Errorf("expected %d, got %d", expected, got)
		}

		userRepository.AssertCalled(t, "CreateUser", user)

	})
}

func userRequest(t testing.TB, username, email, password string) *strings.Reader {
	t.Helper()
	body := `{"username": "` + username + `", "email": "` + email + `", "password": "` + password + `"}`
	return strings.NewReader(body)
}
