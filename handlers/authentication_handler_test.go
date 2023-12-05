package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"ropc-backend/mocks"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FieldValidations(t *testing.T) {

	t.Run("should fail if content type is not application/x-www-form-urlencoded", func(t *testing.T) {
		var (
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			response, request = BuildTestRequest(t, nil)
		)

		assert.Panics(t, func() {
			handler.Authenticate(response, request)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)

	})

	t.Run("should fail with \"client id is required\" if client_id is not provided", func(t *testing.T) {

		var (
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			body              = body("", "secret", "type")
			response, request = BuildTestRequest(t, body)
		)

		assert.PanicsWithError(t, "client id is required", func() {
			handler.Authenticate(response, request)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})

	t.Run("should fail with \"client secret is required\" if client_secret is not provided", func(t *testing.T) {

		var (
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			body              = body("clientId", "", "type")
			response, request = BuildTestRequest(t, body)
		)

		assert.PanicsWithError(t, "client secret is required", func() {
			handler.Authenticate(response, request)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})

	t.Run("should fail with \"grant type is required\" if grant_type is not provided", func(t *testing.T) {

		var (
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			body              = body("clientId", "clientSecret", "")
			response, request = BuildTestRequest(t, body)
		)

		assert.PanicsWithError(t, "grant type is required", func() {
			handler.Authenticate(response, request)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})

	t.Run("should panic with invalid grant type", func(t *testing.T) {

		var (
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			body              = body("clientId", "clientSecret", "invalid_grant_type")
			response, request = BuildTestRequest(t, body)
		)

		assert.PanicsWithError(t, "invalid grant type", func() {
			handler.Authenticate(response, request)
		})

		authenticator.AssertNotCalled(t, "Authenticate", mock.Anything, mock.Anything)
	})
}

func TestAuthenticationHandler_Authenticate(t *testing.T) {

	t.Run("should return 401 with invalid credentials", func(t *testing.T) {

		var (
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			response, request = BuildTestRequest(t, body("clientId", "clientSecret", "client_credentials"))
		)

		authenticator.On("ClientCredentials", "clientId", "clientSecret").Return("", errors.New("auth failed"))
		handler.Authenticate(response, request)

		expected := http.StatusUnauthorized
		got := response.Code

		if expected != got {
			t.Errorf("expected %d but got %d", expected, got)
		}

		authenticator.AssertCalled(t, "ClientCredentials", "clientId", "clientSecret")
	})

	t.Run("should return code 200", func(t *testing.T) {

		var (
			clientId          = "clientId"
			clientSecret      = "clientSecret"
			grantType         = "client_credentials"
			body              = body(clientId, clientSecret, grantType)
			authenticator     = new(mocks.AuthenticatorService)
			ctx               = new(mocks.Context)
			handler           = NewAuthenticationHandler(authenticator, ctx)
			response, request = BuildTestRequest(t, body)
		)

		authenticator.On("ClientCredentials", clientId, clientSecret).Return("", nil)
		handler.Authenticate(response, request)

		expected := http.StatusOK
		got := response.Code

		if expected != got {
			t.Errorf("expected %d but got %d", expected, got)
		}

		authenticator.AssertCalled(t, "ClientCredentials", clientId, clientSecret)

	})
}

func body(clientId, clientSecret, grantType string) io.Reader {

	data := url.Values{}
	data.Set("client_secret", clientSecret)
	//data.Set("username", username)
	//data.Set("password", password)
	data.Set("client_id", clientId)
	data.Set("grant_type", grantType)

	return strings.NewReader(data.Encode())
}
