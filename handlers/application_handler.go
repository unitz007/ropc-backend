package handlers

import (
	"backend-server/model"
	"backend-server/repositories"
	"backend-server/routers"
	"backend-server/utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ApplicationHandler interface {
	CreateApplication(w http.ResponseWriter, r *http.Request)
	GenerateSecret(w http.ResponseWriter, r *http.Request)
	GetApplications(w http.ResponseWriter, r *http.Request)
	GetApplication(w http.ResponseWriter, r *http.Request)
}

type applicationHandler struct {
	router                routers.Router
	applicationRepository repositories.ApplicationRepository
}

func (a *applicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	err, clientId := a.router.GetPathVariable(r, "client_id")
	if err != nil {
		_ = PrintResponse[any](404, w, nil)
	}

	user, _ := utils.GetUserFromContext(r.Context())

	app, err := a.applicationRepository.GetByClientIdAndUserId(clientId, user.ID)
	if err != nil {
		panic(errors.New("application not found"))
	}
	res := model.NewResponse[*model.ApplicationDto]("success", app.ToDTO())

	_ = PrintResponse[*model.Response[*model.ApplicationDto]](http.StatusOK, w, res)

}

func (a *applicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {

	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		panic(err)
	}

	apps := a.applicationRepository.GetAll(user.ID)

	response := make([]*model.ApplicationDto, 0)
	for _, app := range apps {
		r := app.ToDTO()

		response = append(response, r)
	}

	_ = PrintResponse[[]*model.ApplicationDto](http.StatusOK, w, response)
}

func (a *applicationHandler) GenerateSecret(w http.ResponseWriter, r *http.Request) {
	var request *model.CreateApplication

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	user, _ := utils.GetUserFromContext(r.Context())

	app, err := a.applicationRepository.GetByClientIdAndUserId(request.ClientId, user.ID)
	if err != nil {
		panic(err)
	}

	secret := uuid.NewString()
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), 0)

	if err != nil {
		panic(errors.New("could not generate secret"))
	}

	appToUpdate := &model.Application{
		ClientId:     app.ClientId,
		ClientSecret: string(hashed),
	}

	_, err = a.applicationRepository.Update(appToUpdate)

	if err != nil {
		panic(errors.New("could not generate secret"))
	}

	applicationResponse := &model.ApplicationResponse{
		ClientId:     request.ClientId,
		ClientSecret: secret,
		RedirectUrl:  app.RedirectUri,
	}

	response := model.NewResponse[*model.ApplicationResponse]("secret generated successfully", applicationResponse)

	_ = PrintResponse[*model.Response[*model.ApplicationResponse]](http.StatusOK, w, response)

}

func NewApplicationHandler(applicationRepository repositories.ApplicationRepository, router routers.Router) ApplicationHandler {
	return &applicationHandler{
		router:                router,
		applicationRepository: applicationRepository,
	}
}

func (a *applicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {

	var request *model.CreateApplication

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	//if request.RedirectUri == "" {
	//	panic(errors.New("redirect uri is required"))
	//}

	if request.Name == "" {
		panic(errors.New("name is required"))
	}

	alreadyExists, _ := a.applicationRepository.GetByClientId(request.ClientId)
	if alreadyExists != nil {
		panic(errors.New("application with this client id already exists"))
	}

	alreadyExists, _ = a.applicationRepository.GetByName(request.Name)
	if alreadyExists != nil {
		panic(errors.New("application with this name already exists"))
	}

	app := &model.Application{
		ClientId:    uuid.NewString(),
		Name:        request.Name,
		RedirectUri: request.RedirectUri,
	}

	user, _ := utils.GetUserFromContext(r.Context())
	app.User = *user

	err = a.applicationRepository.Create(app)
	if err != nil {
		panic(err)
	}

	response := model.NewResponse[*model.ApplicationDto]("application created successfully", app.ToDTO())

	_ = PrintResponse[*model.Response[*model.ApplicationDto]](http.StatusCreated, w, response)

}
