package handlers

import (
	"backend-server/mocks"
	"backend-server/model"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAppHandler(t *testing.T) {

	var clientHandler ApplicationHandler
	router := new(mocks.Router)
	tt := []struct {
		Name        string
		Body        io.Reader
		ShouldPanic bool
	}{
		{
			Name:        "nil body prepareRequest",
			Body:        nil,
			ShouldPanic: true,
		},
		{
			Name:        "invalid json body",
			Body:        strings.NewReader(`"dhhdhdhr": "ljll"`),
			ShouldPanic: true,
		},
		{
			Name:        "valid json body",
			Body:        strings.NewReader(`{ "client_id": "test_id", "name": "test_name", "redirect_uri": "http://localhost:9090/"}`),
			ShouldPanic: false,
		},
		{
			Name:        "valid json body with empty client_id",
			Body:        strings.NewReader(`{ "client_id": "" }`),
			ShouldPanic: true,
		},
	}

	for _, w := range tt {

		t.Run(w.Name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/clients", w.Body)
			response := httptest.NewRecorder()

			exec := func() {

				repo := new(mocks.ApplicationRepository)
				repo.On("GetByClientId", "test_id").Return(nil, errors.New("application"))
				repo.On("GetByName", "test_name").Return(nil, errors.New("application"))
				repo.On("Create", &model.Application{ClientId: "test_id", Name: "test_name"}).Return(nil)

				clientHandler = NewApplicationHandler(repo, router)
				clientHandler.CreateApplication(response, request)
			}

			if w.ShouldPanic == true {
				assert.Panics(t, exec, "should panic due to invalid prepareRequest body")
			} else {
				assert.NotPanics(t, exec, "should not panic")
			}
		})
	}

	t.Run("should panic with client id is required", func(t *testing.T) {
		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "", "name": "test_secret", "redirect_uri": "http://localhost:9090"}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		appHandler := NewApplicationHandler(applicationRepository, router)
		exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "client id is required", exec)

	})

	t.Run("should panic with name is required", func(t *testing.T) {
		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "jhb", "name": "", "redirect_uri": "http://localhost:8080"}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		appHandler := NewApplicationHandler(applicationRepository, router)
		exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "name is required", exec)
	})

	t.Run("should panic with redirect 'uri is required'", func(t *testing.T) {
		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "jhb", "name": "name", "redirect_uri": ""}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		appHandler := NewApplicationHandler(applicationRepository, router)
		exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "redirect uri is required", exec)
	})

	t.Run("successful prepareRequest should return 201 CREATED", func(t *testing.T) {
		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "test_id", "name": "test_name", "redirect_uri": "http://localhost:3030"}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		applicationRepository.On("GetByClientId", "test_id").Return(nil, errors.New("application"))
		applicationRepository.On("GetByName", "test_name").Return(nil, errors.New("application"))
		applicationRepository.On("Create", &model.Application{ClientId: "test_id", Name: "test_name"}).Return(nil)

		appHandler := NewApplicationHandler(applicationRepository, router)

		appHandler.CreateApplication(response, request)

		expected := http.StatusCreated
		got := response.Code

		if expected != got {
			t.Errorf("expected %v got %v", expected, got)
		}

		applicationRepository.AssertCalled(t, "GetByClientId", "test_id")
		applicationRepository.AssertCalled(t, "GetByName", "test_name")
		applicationRepository.AssertCalled(t, "Create", &model.Application{ClientId: "test_id", Name: "test_name"})

	})

	t.Run("should panic with 'application with this client id already exists'", func(t *testing.T) {

		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "test_id", "name": "test_name", "redirect_uri": "http://localhost:9090"}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		applicationRepository.On("GetByClientId", "test_id").Return(&model.Application{ClientId: "test_id", Name: "test_name"}, nil)

		appHandler := NewApplicationHandler(applicationRepository, router)

		exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "application with this client id already exists", exec)

		applicationRepository.AssertCalled(t, "GetByClientId", "test_id")
		applicationRepository.AssertNotCalled(t, "GetByName", "test_name")
		applicationRepository.AssertNotCalled(t, "Create", &model.Application{ClientId: "test_client", Name: "test_name"})

	})

	t.Run("should panic with 'application with this name already exists'", func(t *testing.T) {

		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "test_id", "name": "test_name", "redirect_uri": "http://localhost:9090/"}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		applicationRepository.On("GetByClientId", "test_id").Return(nil, errors.New("application"))

		applicationRepository.On("GetByName", "test_name").Return(&model.Application{ClientId: "test_id", Name: "test_name"}, nil)

		appHandler := NewApplicationHandler(applicationRepository, router)

		exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "application with this name already exists", exec)

		applicationRepository.AssertCalled(t, "GetByName", "test_name")
		applicationRepository.AssertCalled(t, "GetByClientId", "test_id")
		applicationRepository.AssertNotCalled(t, "Create", &model.Application{ClientId: "test_client", Name: "test_name"})

	})
}

func TestGenerateClientSecret(t *testing.T) {
	router := new(mocks.Router)
	t.Run("should panic with 'application does not exist'", func(t *testing.T) {

		body := strings.NewReader(`{ "client_id": "test_client"}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps/generate-secret", body)
		response := httptest.NewRecorder()

		repoMock := new(mocks.ApplicationRepository)
		repoMock.On("GetByClientId", "test_client").Return(nil, errors.New("application does not exist"))
		repoMock.On("Update", &model.Application{ClientId: "test_client"}).Return(mock.Anything, nil)

		handler := NewApplicationHandler(repoMock, router)

		exec := func() {
			handler.GenerateSecret(response, request)
		}

		assert.PanicsWithError(t, "application does not exist", exec)

		repoMock.AssertCalled(t, "GetByClientId", "test_client")
		repoMock.AssertNotCalled(t, "Update", &model.Application{ClientId: "test_client"})

	})

	t.Run("should generate client secret", func(t *testing.T) {

		t.Skip()

		clientId := "test_client"
		secret := uuid.NewString()
		hashedSecret, _ := bcrypt.GenerateFromPassword([]byte(secret), 0)

		app := &model.Application{ClientId: clientId}

		body := strings.NewReader(`{ "client_id": "test_client"}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps/generate-secret", body)
		response := httptest.NewRecorder()

		repoMock := new(mocks.ApplicationRepository)
		handler := NewApplicationHandler(repoMock, router)

		toUpdate := &model.Application{
			ClientId:     clientId,
			ClientSecret: string(hashedSecret),
		}

		repoMock.On("GetByClientId", clientId).Return(app, nil)
		repoMock.On("Update", toUpdate).Return(toUpdate, nil)

		handler.GenerateSecret(response, request)

		if response.Code != http.StatusOK {
			t.Error("should return 200 OK status code")
		}

		repoMock.AssertCalled(t, "Get", clientId)
		//repoMock.AssertCalled(t, "Update", &model.Application{ClientId: clientId, ClientSecret: string(hashedSecret)})

	})
}

func TestGetApplication(t *testing.T) {

	router := new(mocks.Router)

	t.Run("should return status 404 if {client_id} is not provided", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps", nil)
		response := httptest.NewRecorder()

		repoMock := new(mocks.ApplicationRepository)
		router.On("GetPathVariable", request, "client_id").Return(errors.New(""), "")
		handler := NewApplicationHandler(repoMock, router)

		handler.GetApplication(response, request)
		expected := http.StatusNotFound
		got := response.Code

		if expected != got {
			t.Errorf("got %d, want %d", expected, got)
		}

		repoMock.AssertNotCalled(t, "GetApplication", mock.Anything)
	})

	t.Run("should panic with application not found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps", nil)
		response := httptest.NewRecorder()

		repoMock := new(mocks.ApplicationRepository)
		router.On("GetPathVariable", request, "client_id").Return(errors.New(""), "")
		repoMock.On("GetByClientId", mock.Anything).Return(nil, errors.New("application not found"))
		handler := NewApplicationHandler(repoMock, router)

		exec := func() { handler.GetApplication(response, request) }
		assert.PanicsWithError(t, "application not found", exec)

		repoMock.AssertCalled(t, "GetByClientId", mock.Anything)

	})
}
