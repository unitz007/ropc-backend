package handlers

import (
	"backend-server/mocks"
	"backend-server/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestCreateAppHandler(t *testing.T) {

	var (
		clientHandler ApplicationHandler
		router        = new(mocks.Router)
	)

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
			Name:        "valid json body with empty client_id",
			Body:        strings.NewReader(`{ "client_id": "" }`),
			ShouldPanic: true,
		},
	}

	for _, w := range tt {

		t.Run(w.Name, func(t *testing.T) {

			var (
				app = &model.Application{
					ClientId: uuid.NewString(),
					Name:     "test_name",
				}
				response, request = BuildTestRequest(t, strings.NewReader(`{ "name": "test_name"}`))
				exec              = func() {

					repo := new(mocks.ApplicationRepository)
					repo.On("GetByClientId", "test_id").Return(nil, nil)
					repo.On("GetByName", "test_name").Return(nil, nil)
					repo.On("Create", app).Return(nil)

					clientHandler = NewApplicationHandler(repo, router)
					clientHandler.CreateApplication(response, request)
				}
			)

			if w.ShouldPanic == true {
				assert.Panics(t, exec, "should panic due to invalid prepareRequest body")
			} else {
				assert.NotPanics(t, exec, "should not panic")
			}
		})
	}

	//t.Run("should panic with client id is required", func(t *testing.T) {
	//	t.Skip()
	//	applicationRepository := new(mocks.ApplicationRepository)
	//	body := strings.NewReader(`{ "client_id": "", "name": "test_secret", "redirect_uri": "http://localhost:9090"}`)
	//
	//	request := httptest.NewRequest(http.MethodPost, "/apps", body)
	//	response := httptest.NewRecorder()
	//
	//	appHandler := NewApplicationHandler(applicationRepository, router)
	//	exec := func() { appHandler.CreateApplication(response, request) }
	//
	//	assert.PanicsWithError(t, "client id is required", exec)
	//
	//})

	t.Run("should panic with name is required", func(t *testing.T) {

		var (
			applicationRepository = new(mocks.ApplicationRepository)
			response, request     = BuildTestRequest(t, strings.NewReader(`{ "client_id": "jhb", "name": "", "redirect_uri": "http://localhost:8080"}`))

			appHandler = NewApplicationHandler(applicationRepository, router)
			exec       = func() { appHandler.CreateApplication(response, request) }
		)

		assert.PanicsWithError(t, "name is required", exec)
	})

	t.Run("should panic with redirect 'uri is required'", func(t *testing.T) {
		t.Skip("not a requirement")
		applicationRepository := new(mocks.ApplicationRepository)
		body := strings.NewReader(`{ "client_id": "jhb", "name": "name", "redirect_uri": ""}`)

		request := httptest.NewRequest(http.MethodPost, "/apps", body)
		response := httptest.NewRecorder()

		appHandler := NewApplicationHandler(applicationRepository, router)
		exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "redirect uri is required", exec)
	})

	t.Run("successful prepareRequest should return 201 CREATED", func(t *testing.T) {

		var (
			application           = &model.Application{Name: "test_name"}
			applicationRepository = new(mocks.ApplicationRepository)
			response, request     = BuildTestRequest(t, strings.NewReader(`{ "name": "test_name"}`))
			appHandler            = NewApplicationHandler(applicationRepository, router)
			user, _               = GetUserFromContext(request.Context())
		)

		applicationRepository.On("GetByNameAndUserId", "test_name", user.ID).Return(nil, nil)
		applicationRepository.On("Create", application).Return(nil)

		appHandler.CreateApplication(response, request)

		expected := http.StatusCreated
		got := response.Code

		if expected != got {
			t.Errorf("expected %v got %v", expected, got)
		}

		applicationRepository.AssertCalled(t, "GetByNameAndUserId", "test_name", user.ID)
		applicationRepository.AssertCalled(t, "Create", application)

	})

	t.Run("should panic with 'application with this client id already exists'", func(t *testing.T) {
		t.Skip()

		var (
			router                = new(mocks.Router)
			applicationRepository = new(mocks.ApplicationRepository)
			response, request     = BuildTestRequest(t, strings.NewReader(`{ "client_id": "test_id", "name": "test_name", "redirect_uri": "http://localhost:9090/"}`))
			appHandler            = NewApplicationHandler(applicationRepository, router)
			exec                  = func() { appHandler.CreateApplication(response, request) }
			user, _               = GetUserFromContext(request.Context())
		)

		//applicationRepository := new(mocks.ApplicationRepository)
		//body := strings.NewReader(`{ "name": "test_name"}`)
		//
		//request := httptest.NewRequest(http.MethodPost, "/apps", body)
		//response := httptest.NewRecorder()

		applicationRepository.On("GetByClientIdAndUserId", "test_client_id", user.ID).Return(&model.Application{ClientId: "test_id", Name: "test_name"}, nil)

		//appHandler := NewApplicationHandler(applicationRepository, router)

		//exec := func() { appHandler.CreateApplication(response, request) }

		assert.PanicsWithError(t, "application with this client id already exists", exec)

		applicationRepository.AssertNotCalled(t, "GetByClientIdAndUserId", "test_client_id", user.ID)
		applicationRepository.AssertNotCalled(t, "Create", &model.Application{ClientId: "test_client_id", Name: "test_name"})

	})

	t.Run("should panic with 'application with this name already exists'", func(t *testing.T) {

		var (
			router                = new(mocks.Router)
			applicationRepository = new(mocks.ApplicationRepository)
			response, request     = BuildTestRequest(t, strings.NewReader(`{ "client_id": "test_id", "name": "test_name", "redirect_uri": "http://localhost:9090/"}`))
			appHandler            = NewApplicationHandler(applicationRepository, router)
			exec                  = func() { appHandler.CreateApplication(response, request) }
			user, _               = GetUserFromContext(request.Context())
		)

		applicationRepository.On("GetByNameAndUserId", "test_name", user.ID).Return(&model.Application{ClientId: "test_id", Name: "test_name"}, nil)

		assert.PanicsWithError(t, "application with this name already exists", exec)
		applicationRepository.AssertCalled(t, "GetByNameAndUserId", "test_name", user.ID)
		applicationRepository.AssertNotCalled(t, "Create", &model.Application{ClientId: "test_client", Name: "test_name"})

	})
}

func TestGenerateClientSecret(t *testing.T) {
	router := new(mocks.Router)
	t.Run("should panic with 'application does not exist'", func(t *testing.T) {

		body := strings.NewReader(`{ "client_id": "test_client"}`)

		request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps/generate-secret", body)
		request = request.WithContext(context.WithValue(request.Context(), UserKey, &model.User{Model: gorm.Model{ID: uint(2)}}))

		response := httptest.NewRecorder()

		repoMock := new(mocks.ApplicationRepository)
		repoMock.On("GetByClientIdAndUserId", "test_client", uint(2)).Return(nil, errors.New("application does not exist"))
		repoMock.On("Update", &model.Application{Name: "test_name"}).Return(mock.Anything, nil)

		handler := NewApplicationHandler(repoMock, router)

		exec := func() {
			handler.GenerateSecret(response, request)
		}

		assert.PanicsWithError(t, "application does not exist", exec)

		repoMock.AssertCalled(t, "GetByClientIdAndUserId", "test_client", uint(2))
		repoMock.AssertNotCalled(t, "Update", &model.Application{Name: "test_name"})

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

		repoMock.On("GetByClientIdAndUserId", clientId).Return(app, nil)
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

	t.Run("should return status 404 if {client_id} is not provided", func(t *testing.T) {

		var (
			response, request = BuildTestRequest(t, nil)
			repoMock          = new(mocks.ApplicationRepository)
			router            = new(mocks.Router)
			handler           = NewApplicationHandler(repoMock, router)
		)

		router.On("GetPathVariable", request, "client_id").Return(errors.New(""), "")

		handler.GetApplication(response, request)
		expected := http.StatusNotFound
		got := response.Code

		if expected != got {
			t.Errorf("got %d, want %d", expected, got)
		}

		repoMock.AssertNotCalled(t, "GetApplication", mock.Anything)
	})

	t.Run("should panic with application not found", func(t *testing.T) {

		var (
			response, request = BuildTestRequest(t, nil)
			router            = new(mocks.Router)
			repoMock          = new(mocks.ApplicationRepository)
			handler           = NewApplicationHandler(repoMock, router)
			exec              = func() { handler.GetApplication(response, request) }
		)

		router.On("GetPathVariable", request, "client_id").Return(nil, "2")
		repoMock.On("GetByClientIdAndUserId", mock.Anything, uint(0)).Return(nil, errors.New("application not found"))

		assert.PanicsWithError(t, "application not found", exec)
		repoMock.AssertCalled(t, "GetByClientIdAndUserId", mock.Anything, uint(0))

	})
}

func TestDeleteApplication(t *testing.T) {

	testApp := model.Application{ClientId: "test_client_id", Model: gorm.Model{ID: uint(2)}}

	t.Run("should panic with application does not exist", func(t *testing.T) {

		var (
			applicationRepository = new(mocks.ApplicationRepository)
			router                = new(mocks.Router)
			response, request     = BuildTestRequest(t, nil)
			handler               = NewApplicationHandler(applicationRepository, router)
			exec                  = func() { handler.DeleteApplication(response, request) }
		)

		router.On("GetPathVariable", request, "client_id").Return(nil, "test_client_id")
		applicationRepository.On("GetByClientIdAndUserId", "test_client_id", testApp.User.ID).Return(nil, errors.New("application does not exist"))

		assert.PanicsWithError(t, "application does not exist", exec)
		router.AssertCalled(t, "GetPathVariable", request, "client_id")
		applicationRepository.AssertCalled(t, "GetByClientIdAndUserId", "test_client_id", testApp.UserID)
		applicationRepository.AssertNotCalled(t, "Delete", mock.Anything)
	})

	t.Run("should return 200 if deleted successfully", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.ApplicationRepository)
			router                = new(mocks.Router)
			response, request     = BuildTestRequest(t, nil)
			handler               = NewApplicationHandler(applicationRepository, router)
		)

		router.On("GetPathVariable", request, "client_id").Return(nil, "test_client_id")
		applicationRepository.On("GetByClientIdAndUserId", "test_client_id", testApp.UserID).Return(&testApp, nil)
		applicationRepository.On("Delete", testApp.ID).Return(nil)

		handler.DeleteApplication(response, request)

		assert.Equal(t, http.StatusOK, response.Code, fmt.Sprintf("expected %d but got %d", http.StatusOK, response.Code))

		router.AssertCalled(t, "GetPathVariable", request, "client_id")
		applicationRepository.AssertCalled(t, "GetByClientIdAndUserId", "test_client_id", testApp.UserID)
		applicationRepository.AssertCalled(t, "Delete", testApp.ID)
	})

	t.Run("should panic with 'failed to delete application' if delete application fails", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.ApplicationRepository)
			router                = new(mocks.Router)
			response, request     = BuildTestRequest(t, nil)
			handler               = NewApplicationHandler(applicationRepository, router)
			exec                  = func() { handler.DeleteApplication(response, request) }
		)

		router.On("GetPathVariable", request, "client_id").Return(nil, "test_client_id")
		applicationRepository.On("GetByClientIdAndUserId", "test_client_id", testApp.UserID).Return(&testApp, nil)
		applicationRepository.On("Delete", testApp.ID).Return(errors.New("dummy errors"))

		assert.PanicsWithError(t, "failed to delete application", exec)
		router.AssertCalled(t, "GetPathVariable", request, "client_id")
		applicationRepository.AssertCalled(t, "GetByClientIdAndUserId", "test_client_id", testApp.UserID)
		applicationRepository.AssertCalled(t, "Delete", testApp.ID)
	})

}

func TestGetApplications(t *testing.T) {
	var (
		applicationRepository = new(mocks.ApplicationRepository)
		router                = new(mocks.Router)
		response, request     = BuildTestRequest(t, nil)
		handler               = NewApplicationHandler(applicationRepository, router)
		applications          = []model.Application{
			{ClientId: "test_id"},
		}
		expectedMessage = fmt.Sprintf("%d application(s) fetched successfully", len(applications))
		user, _         = GetUserFromContext(request.Context())
	)

	applicationRepository.On("GetAll", user.ID).Return(applications)
	handler.GetApplications(response, request)

	var responseBody *model.Response[[]*model.Application]

	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal("could  not unmarshal")
	}

	assert.Equal(t, expectedMessage, responseBody.Message)
}
