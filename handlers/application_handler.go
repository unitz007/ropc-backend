package handlers

import (
	"backend-server/model"
	"backend-server/repositories"
	"errors"
	"net/http"
	"strings"

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
	applicationRepository repositories.ApplicationRepository
}

func (a *applicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")

	isHtml := strings.Contains(accept, "text/html")

	if isHtml == true {
		http.ServeFile(w, r, "./assets/html/create_application.html")
		return
	} else {
		panic(errors.New("accept type not supported"))
	}
}

func (a *applicationHandler) GetApplications(w http.ResponseWriter, _ *http.Request) {
	apps := a.applicationRepository.GetAll()
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

	app, err := a.applicationRepository.GetByClientId(request.ClientId)
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
	}

	_ = PrintResponse[*model.ApplicationResponse](http.StatusOK, w, applicationResponse)

}

func NewApplicationHandler(applicationRepository repositories.ApplicationRepository) ApplicationHandler {
	return &applicationHandler{
		applicationRepository: applicationRepository,
	}
}

func (a *applicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {

	var request *model.CreateApplication

	err := JsonToStruct(r.Body, &request)
	if err != nil {
		panic(errors.New("invalid request body"))
	}

	if request.ClientId == "" {
		panic(errors.New("client id is required"))
	}

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

	app, err := model.NewApplication(request.ClientId, request.Name)
	if err != nil {
		panic(err)
	}

	err = a.applicationRepository.Create(app)
	if err != nil {
		panic(err)
	}

	response := model.NewResponse[*model.ApplicationDto]("application created successfully", app.ToDTO())

	_ = PrintResponse[*model.Response[*model.ApplicationDto]](http.StatusCreated, w, response)

}
