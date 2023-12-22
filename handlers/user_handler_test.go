package handlers

import (
	"net/http"
	"net/http/httptest"
	"ropc-backend/kernel"
	"ropc-backend/mocks"
	"ropc-backend/model"
	"ropc-backend/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	username := "username"
	email := "test@gmail.com"
	password := "password"

	user := model.User{
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
		userRepository := new(mocks.Repository[model.User])
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
		userRepository := new(mocks.Repository[model.User])
		config := new(mocks.Config)
		handler := NewUserHandler(config, userRepository)

		json := userRequest(t, username, email, password)

		userRepository.On("Create", user).Return(kernel.EntityAlreadyExists)

		request := httptest.NewRequest(http.MethodPost, "/users", json)
		response := httptest.NewRecorder()

		exec := func() {
			handler.CreateUser(response, request)
		}

		assert.Panics(t, exec)
		userRepository.AssertCalled(t, "Create", user)
	})

	t.Run("user created successfully", func(t *testing.T) {
		userRepository := new(mocks.Repository[model.User])
		json := userRequest(t, username, email, password)

		userRepository.On("Create", user).Return(nil)

		request := httptest.NewRequest(http.MethodPost, "/users", json)
		response := httptest.NewRecorder()

		handler := NewUserHandler(utils.NewConfig(), userRepository)
		handler.CreateUser(response, request)

		expected := http.StatusCreated
		got := response.Code

		if expected != got {
			t.Errorf("expected %d, got %d", expected, got)
		}

		userRepository.AssertCalled(t, "Create", user)

	})
}

func userRequest(t testing.TB, username, email, password string) *strings.Reader {

	t.Helper()
	body := `{"username": "` + username + `", "email": "` + email + `", "password": "` + password + `"}`
	return strings.NewReader(body)
}
