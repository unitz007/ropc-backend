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
	"gorm.io/gorm"
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

	user, err := GetUserFromContext(request.Context())
	if err != nil {
		_ = utils.PrintResponse[any](http.StatusForbidden, response, nil)
		return
	}

	condition := utils.Queries[utils.WhereClientIdAndUserIdIs](clientId, user.ID)

	app, err := a.repository.Get(condition)

	if err != nil {
		panic(err)
	}

	condition = utils.Queries[utils.WhereIdIs](app.ID)

	err = a.repository.Delete(condition)
	if err != nil {
		panic(errors.New("failed to delete application"))
	}

	body := *model.NewResponse[any]("application deleted successfully", nil)

	_ = utils.PrintResponse[model.Response[any]](http.StatusOK, response, body)
}

func (a *applicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	err, clientId := a.Router().GetPathVariable(r, "client_id")
	if err != nil {
		_ = utils.PrintResponse[any](404, w, nil)
		return
	}

	user, _ := GetUserFromContext(r.Context())

	condition := utils.Queries[utils.WhereClientIdAndUserIdIs](clientId, user.ID)

	app, err := a.repository.Get(condition)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			panic(errors.New("application not found"))
		default:
			panic(err.Error())
		}

	}
	res := model.NewResponse[*model.ApplicationDto]("success", app.ToDTO())

	_ = utils.PrintResponse[*model.Response[*model.ApplicationDto]](http.StatusOK, w, res)

}

func (a *applicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {

	user, err := GetUserFromContext(r.Context())
	if err != nil {
		panic(err)
	}

	condition := utils.Queries[utils.WhereUserIdIs](user.ID)

	apps := a.repository.GetAll(condition)

	response := make([]*model.ApplicationDto, 0)
	for _, app := range apps {
		r := app.ToDTO()
		response = append(response, r)
	}

	responseBody := model.NewResponse[[]*model.ApplicationDto](fmt.Sprintf("%d application(s) fetched successfully", len(apps)), response)

	_ = utils.PrintResponse[*model.Response[[]*model.ApplicationDto]](http.StatusOK, w, responseBody)
}

func (a *applicationHandler) GenerateSecret(w http.ResponseWriter, r *http.Request) {
	var request *model.CreateApplication

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	user, err := GetUserFromContext(r.Context())
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

	response := model.NewResponse[*model.ApplicationResponse]("secret generated successfully", applicationResponse)

	_ = utils.PrintResponse[*model.Response[*model.ApplicationResponse]](http.StatusOK, w, response)

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
		user, _ = GetUserFromContext(r.Context())
	)

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	if request.Name == utils.Blank {
		panic(errors.New("name is required"))
	}

	condition := utils.Queries[utils.WhereNameAndUserIdIs](request.Name, user.ID)

	alreadyExists, _ := a.repository.Get(condition)
	if alreadyExists != nil {
		panic(errors.New("application with this name already exists"))
	}

	app := model.Application{
		Name:        request.Name,
		RedirectUri: request.RedirectUri,
		User:        *user,
		ClientId:    uuid.NewString(),
	}

	err = a.repository.Create(app)
	if err != nil {
		panic(err)
	}

	response := model.NewResponse[any]("application created successfully", nil)

	_ = utils.PrintResponse[*model.Response[any]](http.StatusCreated, w, response)
}
