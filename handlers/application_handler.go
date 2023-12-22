package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	applicationNotFoundErr = "Application does not exist"
)

type ApplicationHandler interface {
	CreateApplication(w http.ResponseWriter, r *http.Request)
	GenerateSecret(w http.ResponseWriter, r *http.Request)
	GetApplications(w http.ResponseWriter, r *http.Request)
	GetApplication(w http.ResponseWriter, r *http.Request)
	DeleteApplication(response http.ResponseWriter, request *http.Request)
}

type applicationHandler struct {
	kernel.Context
	repository kernel.Repository[model.Application]
}

func (a *applicationHandler) DeleteApplication(response http.ResponseWriter, request *http.Request) {

	err, clientId := a.Router().GetPathVariable(request, "client_id")
	if err != nil {
		panic(err)
	}

	user := GetUserFromContext(request.Context())

	condition := utils.Queries[utils.WhereClientIdAndUserIdIs](clientId, user.ID)

	app, err := a.repository.Get(condition)

	if err != nil {
		if errors.Is(err, kernel.EntityNotFoundError) {
			_ = utils.PrintResponseNew[any](response, http.StatusNotFound, applicationNotFoundErr, nil)
			return
		} else {
			panic(err)
		}
	}

	condition = utils.Queries[utils.WhereIdIs](app.ID)

	err = a.repository.Delete(condition)
	if err != nil {
		panic(errors.New("failed to delete application"))
	}

	_ = utils.PrintResponseNew[any](response, http.StatusOK, "Application deleted successfully", nil)
}

func (a *applicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	err, clientId := a.Router().GetPathVariable(r, "client_id")
	if err != nil {
		_ = utils.PrintResponse[any](404, w, nil)
		return
	}

	user := GetUserFromContext(r.Context())

	condition := utils.Queries[utils.WhereClientIdAndUserIdIs](clientId, user.ID)

	app, err := a.repository.Get(condition)

	if err != nil {
		if errors.Is(err, kernel.EntityNotFoundError) {
			panic(applicationNotFoundErr)
		} else {
			panic(err.Error())
		}
	}

	_ = utils.PrintResponseNew[*model.ApplicationDto](w, http.StatusOK, "Application fetched successfully", app.ToDTO())

}

func (a *applicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {

	user := GetUserFromContext(r.Context())

	condition := utils.Queries[utils.WhereUserIdIs](user.ID)

	apps := a.repository.GetAll(condition)

	appDtos := make([]*model.ApplicationDto, 0)
	for _, app := range apps {
		r := app.ToDTO()
		appDtos = append(appDtos, r)
	}

	_ = utils.PrintResponseNew[[]*model.ApplicationDto](w, http.StatusOK, fmt.Sprintf("%d application(s) fetched successfully", len(apps)), appDtos)
}

func (a *applicationHandler) GenerateSecret(w http.ResponseWriter, r *http.Request) {
	var request *model.CreateApplication

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	user := GetUserFromContext(r.Context())
	if err != nil {
		panic(errors.New("forbidden: not allowed to make this request"))
	}

	condition := utils.Queries[utils.WhereClientIdAndUserIdIs](request.ClientId, user.ID)
	app, err := a.repository.Get(condition)
	if err != nil {
		panic(err)
	}

	secret := uuid.NewString()
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), 0)

	if err != nil {
		panic(errors.New("could not generate secret"))
	}

	update := map[string]any{
		"ClientSecret": string(hashed),
	}

	err = a.repository.Update(app.ID, update)

	if err != nil {
		panic(err)
	}

	applicationResponse := &model.ApplicationResponse{
		ClientId:     request.ClientId,
		ClientSecret: secret,
		RedirectUrl:  app.RedirectUri,
	}

	_ = utils.PrintResponseNew[*model.ApplicationResponse](w, http.StatusOK, "Secret generated successfully", applicationResponse)

}
func NewApplicationHandler(applicationRepository kernel.Repository[model.Application], ctx kernel.Context) ApplicationHandler {
	return &applicationHandler{
		Context:    ctx,
		repository: applicationRepository,
	}
}

func (a *applicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {

	var (
		request *model.CreateApplication
		user    = GetUserFromContext(r.Context())
	)

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	if request.Name == utils.Blank {
		panic(errors.New("name is required"))
	}

	if request.ClientId == utils.Blank {
		panic(errors.New("client_id is required"))
	}

	clientId, err := uuid.Parse(request.ClientId)
	if err != nil {
		panic(errors.New("invalid client_id. client_id should be uuid"))
	}

	condition := utils.Queries[utils.WhereNameAndUserIdIs](request.Name, user.ID)

	alreadyExists, _ := a.repository.Get(condition)
	if alreadyExists != nil {
		panic(kernel.NewError(http.StatusConflict, "Oops, you have an application with this name already."))
	}

	condition = utils.Queries[utils.WhereClientIdAndUserIdIs](request.ClientId, user.ID)
	alreadyExists, _ = a.repository.Get(condition)
	if alreadyExists != nil {
		panic(kernel.NewError(http.StatusConflict, "Oops, you have an application with this client id already."))
	}

	app := model.Application{
		Name:        request.Name,
		RedirectUri: request.RedirectUri,
		User:        *user,
		ClientId:    clientId.String(),
	}

	err = a.repository.Create(app)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Location", "/apps/"+app.ClientId)

	_ = utils.PrintResponseNew[any](w, http.StatusCreated, "Application created successfully", nil)
}
