package handlers

import (
	"backend-server/mocks"
	"backend-server/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FieldValidations(t *testing.T) {

	t.Run("should fail if content type is not application/x-www-form-urlencoded", func(t *testing.T) {
		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)

		req := httptest.NewRequest(http.MethodPost, "/tokens", nil)
		w := httptest.NewRecorder()
		assert.Panics(t, func() {
			handler.Authenticate(w, req)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)

	})

	t.Run("should fail with \"client id is required\" if client_id is not provided", func(t *testing.T) {
		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)
		req := prepareRequest("", "secret", "type")
		w := httptest.NewRecorder()
		assert.PanicsWithError(t, "client id is required", func() {
			handler.Authenticate(w, req)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})

	t.Run("should fail with \"client secret is required\" if client_secret is not provided", func(t *testing.T) {
		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)
		req := prepareRequest("clientId", "", "type")
		w := httptest.NewRecorder()
		assert.PanicsWithError(t, "client secret is required", func() {
			handler.Authenticate(w, req)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})

	t.Run("should fail with \"grant type is required\" if grant_type is not provided", func(t *testing.T) {
		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)
		req := prepareRequest("clientId", "clientSecret", "")
		w := httptest.NewRecorder()
		assert.PanicsWithError(t, "grant type is required", func() {
			handler.Authenticate(w, req)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})

	t.Run("should panic with invalid grant type", func(t *testing.T) {
		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)
		req := prepareRequest("clientId", "clientSecret", "invalid_grant_type")
		w := httptest.NewRecorder()
		assert.PanicsWithError(t, "invalid grant type", func() {
			handler.Authenticate(w, req)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})
}

func TestAuthenticationHandler_Authenticate(t *testing.T) {

	clientId := "clientId"
	clientSecret := "clientSecret"
	grantType := "client_credentials"

	req := prepareRequest(clientId, clientSecret, grantType)

	t.Run("should return 401 with invalid credentials", func(t *testing.T) {

		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)
		w := httptest.NewRecorder()

		authenticator.On("ClientCredentials", clientId, clientSecret).Return(nil, errors.New("auth failed"))
		handler.Authenticate(w, req)

		expected := http.StatusUnauthorized
		got := w.Code

		if expected != got {
			t.Errorf("expected %d but got %d", expected, got)
		}

		authenticator.AssertCalled(t, "ClientCredentials", clientId, clientSecret)
	})

	t.Run("should return code 200", func(t *testing.T) {
		authenticator := new(mocks.AuthenticatorService)
		handler := NewAuthenticationHandler(authenticator)
		w := httptest.NewRecorder()

		authenticator.On("ClientCredentials", clientId, clientSecret).Return(&model.Token{}, nil)
		handler.Authenticate(w, req)

		expected := http.StatusOK
		got := w.Code

		if expected != got {
			t.Errorf("expected %d but got %d", expected, got)
		}

		authenticator.AssertCalled(t, "ClientCredentials", clientId, clientSecret)

	})
}

func prepareRequest(clientId, clientSecret, grantType string) *http.Request {

	data := url.Values{}
	data.Set("client_secret", clientSecret)
	//data.Set("username", username)
	//data.Set("password", password)
	data.Set("client_id", clientId)
	data.Set("grant_type", grantType)

	req := httptest.NewRequest(http.MethodPost, "/tokens", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}
