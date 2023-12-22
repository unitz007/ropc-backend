package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"ropc-backend/kernel"
	"ropc-backend/mocks"
	"ropc-backend/model"
	"ropc-backend/utils"
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
		ctx           = new(mocks.Context)
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

					repo := new(mocks.Repository[model.Application])
					repo.On("GetByClientId", "test_id").Return(nil, nil)
					repo.On("GetByName", "test_name").Return(nil, nil)
					repo.On("Create", app).Return(nil)

					clientHandler = NewApplicationHandler(repo, ctx)
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

	t.Run("should fail when name is not provided", func(t *testing.T) {

		var (
			applicationRepository = new(mocks.Repository[model.Application])
			response, request     = BuildTestRequest(t, strings.NewReader(`{ "client_id": "jhb", "name": "", "redirect_uri": "http://localhost:8080"}`))
			ctx                   = new(mocks.Context)
			appHandler            = NewApplicationHandler(applicationRepository, ctx)
			exec                  = func() { appHandler.CreateApplication(response, request) }
		)

		assert.Panics(t, exec)

		applicationRepository.AssertNotCalled(t, "Get", utils.Queries[utils.WhereUserIdIs], mock.Anything)
		applicationRepository.AssertNotCalled(t, "Create", model.Application{})
	})

	t.Run("should panic when client_id is not provided", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.Repository[model.Application])
			response, request     = BuildTestRequest(t, strings.NewReader(`{"name": "name", "redirect_uri": ""}`))
			ctx                   = new(mocks.Context)
			appHandler            = NewApplicationHandler(applicationRepository, ctx)
			exec                  = func() { appHandler.CreateApplication(response, request) }
		)

		assert.Panics(t, exec)

		applicationRepository.AssertNotCalled(t, "Get", utils.Queries[utils.WhereUserIdIs], mock.Anything)
		applicationRepository.AssertNotCalled(t, "Create", model.Application{})

	})

	t.Run("should panic when client_id is not uuid", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.Repository[model.Application])
			response, request     = BuildTestRequest(t, strings.NewReader(`{"name": "name", "redirect_uri": "", "client_id": "uyuuhiuhi"}`))
			ctx                   = new(mocks.Context)
			appHandler            = NewApplicationHandler(applicationRepository, ctx)
			exec                  = func() { appHandler.CreateApplication(response, request) }
			user                  = GetUserFromContext(request.Context())
		)

		assert.Panics(t, exec)

		applicationRepository.AssertNotCalled(t, "Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID))
		applicationRepository.AssertNotCalled(t, "Create", model.Application{})

	})

	t.Run("should fail if client_id already exists'", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.Repository[model.Application])
			response, request     = BuildTestRequest(t, strings.NewReader(`{"name": "name", "redirect_uri": "", "client_id": "f020a2d9-f07d-4ace-bdb6-c2d17332a10d"}`))
			ctx                   = new(mocks.Context)
			appHandler            = NewApplicationHandler(applicationRepository, ctx)
			exec                  = func() { appHandler.CreateApplication(response, request) }
			user                  = GetUserFromContext(request.Context())
			application           = &model.Application{Name: "test_name", User: *user}
		)

		applicationRepository.On("Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID)).Return(nil, errors.New(""))
		applicationRepository.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("f020a2d9-f07d-4ace-bdb6-c2d17332a10d", user.ID)).Return(application, nil)

		assert.Panics(t, exec)

		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID))
		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("f020a2d9-f07d-4ace-bdb6-c2d17332a10d", user.ID))
		applicationRepository.AssertNotCalled(t, "Create", application)

	})

	t.Run("should be successful", func(t *testing.T) {

		var (
			applicationRepository = new(mocks.Repository[model.Application])
			response, request     = BuildTestRequest(t, strings.NewReader(`{ "name": "name", "client_id":"f020a2d9-f07d-4ace-bdb6-c2d17332a10d"}`))
			ctx                   = new(mocks.Context)
			appHandler            = NewApplicationHandler(applicationRepository, ctx)
			user                  = GetUserFromContext(request.Context())
			application           = model.Application{Name: "name", ClientId: "f020a2d9-f07d-4ace-bdb6-c2d17332a10d", User: *user}
		)

		applicationRepository.On("Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID)).Return(nil, errors.New(""))
		applicationRepository.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("f020a2d9-f07d-4ace-bdb6-c2d17332a10d", user.ID)).Return(nil, errors.New(""))
		applicationRepository.On("Create", application).Return(nil)

		appHandler.CreateApplication(response, request)

		expected := http.StatusCreated
		got := response.Code

		if expected != got {
			t.Errorf("expected %v got %v", expected, got)
		}

		contentLocation := response.Header().Get("Content-Location")
		if contentLocation == utils.Blank {
			t.Errorf("expected %v got %v", "/apps/"+application.ClientId, contentLocation)
		}

		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID))
		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("f020a2d9-f07d-4ace-bdb6-c2d17332a10d", user.ID))
		applicationRepository.AssertCalled(t, "Create", application)

	})

	t.Run("should panic when application name already exists", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.Repository[model.Application])
			response, request     = BuildTestRequest(t, strings.NewReader(`{"name": "name", "redirect_uri": "", "client_id": "f020a2d9-f07d-4ace-bdb6-c2d17332a10d"}`))
			ctx                   = new(mocks.Context)
			appHandler            = NewApplicationHandler(applicationRepository, ctx)
			exec                  = func() { appHandler.CreateApplication(response, request) }
			user                  = GetUserFromContext(request.Context())
			application           = &model.Application{Name: "name", User: *user}
		)

		applicationRepository.On("Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID)).Return(application, nil)

		assert.Panics(t, exec)

		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereNameAndUserIdIs]("name", user.ID))
		applicationRepository.AssertNotCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("f020a2d9-f07d-4ace-bdb6-c2d17332a10d", user.ID))
		applicationRepository.AssertNotCalled(t, "Create", application)

	})

}

func TestGenerateClientSecret(t *testing.T) {
	ctx := new(mocks.Context)
	t.Run("should panic when 'application does not exist'", func(t *testing.T) {

		var (
			response, request = BuildTestRequest(t, strings.NewReader(`{ "client_id": "test_client"}`))
			repoMock          = new(mocks.Repository[model.Application])
			handler           = NewApplicationHandler(repoMock, ctx)
			user              = GetUserFromContext(request.Context())
			exec              = func() { handler.GenerateSecret(response, request) }
		)

		repoMock.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client", user.ID)).Return(nil, errors.New("application does not exist"))
		repoMock.On("Update", &model.Application{Name: "test_name"}).Return(mock.Anything, nil)

		assert.Panics(t, exec)

		repoMock.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client", user.ID))
		repoMock.AssertNotCalled(t, "Update", &model.Application{Name: "test_name"})

	})

	t.Run("should generate client secret", func(t *testing.T) {
		t.Skip()

		var (
			clientId          = "test_client"
			secret            = uuid.NewString()
			hashedSecret, _   = bcrypt.GenerateFromPassword([]byte(secret), 0)
			app               = &model.Application{ClientId: clientId}
			response, request = BuildTestRequest(t, strings.NewReader(`{ "client_id": "test_client"}`))
			repoMock          = new(mocks.Repository[model.Application])
			ctx               = new(mocks.Context)
			user              = GetUserFromContext(request.Context())
			handler           = NewApplicationHandler(repoMock, ctx)
		)

		update := map[string]any{
			"ClientSecret": string(hashedSecret),
		}

		repoMock.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client", user.ID)).Return(app, nil)
		repoMock.On("Update", app.ID, update).Return(nil)

		handler.GenerateSecret(response, request)

		if response.Code != http.StatusOK {
			t.Error("should return 200 OK status code")
		}

		//repoMock.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client", user.ID))
		//repoMock.AssertCalled(t, "Update", toUpdate)

	})
}

func TestGetApplication(t *testing.T) {

	t.Run("should return status 404 if {client_id} is not provided", func(t *testing.T) {

		var (
			response, request = BuildTestRequest(t, nil)
			repoMock          = new(mocks.Repository[model.Application])
			ctx               = new(mocks.Context)
			handler           = NewApplicationHandler(repoMock, ctx)
			router            = new(mocks.Router)
		)

		ctx.On("Router").Return(router)
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
			ctx               = new(mocks.Context)
			repoMock          = new(mocks.Repository[model.Application])
			handler           = NewApplicationHandler(repoMock, ctx)
			exec              = func() { handler.GetApplication(response, request) }
			router            = new(mocks.Router)
			user              = GetUserFromContext(request.Context())
		)

		ctx.On("Router").Return(router)
		router.On("GetPathVariable", request, "client_id").Return(nil, "2")
		repoMock.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("2", user.ID)).Return(nil, kernel.EntityNotFoundError)

		assert.Panics(t, exec)
		repoMock.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("2", user.ID))

	})
}

func TestDeleteApplication(t *testing.T) {

	t.Run("should fail when application does not exist", func(t *testing.T) {

		var (
			applicationRepository = new(mocks.Repository[model.Application])
			ctx                   = new(mocks.Context)
			response, request     = BuildTestRequest(t, nil)
			handler               = NewApplicationHandler(applicationRepository, ctx)
			user                  = GetUserFromContext(request.Context())
			testApp               = model.Application{ClientId: "test_client_id", User: *user, Model: gorm.Model{ID: uint(1)}}
			router                = new(mocks.Router)
		)

		ctx.On("Router").Return(router)
		router.On("GetPathVariable", request, "client_id").Return(nil, "test_client_id")
		applicationRepository.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client_id", testApp.User.ID)).Return(nil, kernel.EntityNotFoundError)

		handler.DeleteApplication(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("expected code %d, got %d", http.StatusNotFound, response.Code)
		}

		var responseBody model.Response[*model.Application]

		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		if err != nil {
			t.Fatal("could  not unmarshal")
		}

		if responseBody.Message != "Application does not exist" {
			t.Errorf("expected %s but got %s", "Application does not exist", responseBody.Message)
		}

		ctx.AssertCalled(t, "Router")
		router.AssertCalled(t, "GetPathVariable", request, "client_id")
		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client_id", testApp.User.ID))
		applicationRepository.AssertNotCalled(t, "Delete", mock.Anything)
	})

	t.Run("should return 200 if deleted successfully", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.Repository[model.Application])
			ctx                   = new(mocks.Context)
			response, request     = BuildTestRequest(t, nil)
			handler               = NewApplicationHandler(applicationRepository, ctx)
			router                = new(mocks.Router)
			user                  = GetUserFromContext(request.Context())
			testApp               = model.Application{Model: gorm.Model{ID: uint(1)}, ClientId: "test_client_id", User: *user}
		)

		ctx.On("Router").Return(router)
		router.On("GetPathVariable", request, "client_id").Return(nil, "test_client_id")

		applicationRepository.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client_id", testApp.User.ID)).Return(&testApp, nil)
		applicationRepository.On("Delete", utils.Queries[utils.WhereIdIs](testApp.ID)).Return(nil)

		handler.DeleteApplication(response, request)

		assert.Equal(t, http.StatusOK, response.Code, fmt.Sprintf("expected %d but got %d", http.StatusOK, response.Code))

		ctx.AssertCalled(t, "Router")
		router.AssertCalled(t, "GetPathVariable", request, "client_id")
		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client_id", testApp.User.ID))
		applicationRepository.AssertCalled(t, "Delete", utils.Queries[utils.WhereIdIs](testApp.ID))

	})

	t.Run("should panic with 'failed to delete application' if delete application fails", func(t *testing.T) {
		var (
			applicationRepository = new(mocks.Repository[model.Application])
			ctx                   = new(mocks.Context)
			response, request     = BuildTestRequest(t, nil)
			handler               = NewApplicationHandler(applicationRepository, ctx)
			exec                  = func() { handler.DeleteApplication(response, request) }
			user                  = GetUserFromContext(request.Context())
			testApp               = model.Application{ClientId: "test_client_id", User: *user, Model: gorm.Model{ID: uint(1)}}
			router                = new(mocks.Router)
		)

		ctx.On("Router").Return(router)
		router.On("GetPathVariable", request, "client_id").Return(nil, "test_client_id")

		applicationRepository.On("Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client_id", testApp.User.ID)).Return(&testApp, nil)
		applicationRepository.On("Delete", utils.Queries[utils.WhereIdIs](testApp.ID)).Return(errors.New("error deleting"))

		assert.PanicsWithError(t, "failed to delete application", exec)

		ctx.AssertCalled(t, "Router")
		router.AssertCalled(t, "GetPathVariable", request, "client_id")
		applicationRepository.AssertCalled(t, "Get", utils.Queries[utils.WhereClientIdAndUserIdIs]("test_client_id", testApp.User.ID))
		applicationRepository.AssertCalled(t, "Delete", utils.Queries[utils.WhereIdIs](testApp.ID))
	})

}

func TestGetApplications(t *testing.T) {
	var (
		applicationRepository = new(mocks.Repository[model.Application])
		ctx                   = new(mocks.Context)
		response, request     = BuildTestRequest(t, nil)
		handler               = NewApplicationHandler(applicationRepository, ctx)
		applications          = []model.Application{
			{ClientId: "test_id"},
		}
		expectedMessage = fmt.Sprintf("%d application(s) fetched successfully", len(applications))
		user            = GetUserFromContext(request.Context())
	)

	condition := utils.Queries[utils.WhereUserIdIs](user.ID)

	applicationRepository.On("GetAll", condition).Return(applications)
	handler.GetApplications(response, request)

	var responseBody *model.Response[[]*model.Application]

	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal("could  not unmarshal")
	}

	assert.Equal(t, expectedMessage, responseBody.Message)
}
