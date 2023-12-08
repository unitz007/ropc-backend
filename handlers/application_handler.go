package handlers

import (
	"errors"
	"net/http"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/utils"
)

type ApplicationHandler interface {
	// DeleteApplication CreateApplication(w http.ResponseWriter, r *http.Request)
	//GenerateSecret(w http.ResponseWriter, r *http.Request)
	//GetApplications(w http.ResponseWriter, r *http.Request)
	//GetApplication(w http.ResponseWriter, r *http.Request)
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
		panic(errors.New("application does not exist"))
	}

	qModel := model.Application{
		UserID:   user.ID,
		ClientId: clientId,
	}

	app, err := a.repository.QuerySingleResult(qModel)

	if err != nil {
		panic(err)
	}

	err = a.repository.Delete(app.ID)
	if err != nil {
		panic(errors.New("failed to delete application"))
	}

	body := *model.NewResponse[any]("application deleted successfully", nil)

	_ = utils.PrintResponse[model.Response[any]](http.StatusOK, response, body)
}

//	NewApplicationHandler func (a *applicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
//		err, clientId := a.Router().GetPathVariable(r, "client_id")
//		if err != nil {
//			_ = utils.PrintResponse[any](404, w, nil)
//			return
//		}
//
//		user, _ := GetUserFromContext(r.Context())
//
//		app, err := a.applicationRepository.GetByClientIdAndUserId(clientId, user.ID)
//		if err != nil {
//			panic(errors.New("application not found"))
//		}
//		res := model.NewResponse[*model.ApplicationDto]("success", app.ToDTO())
//
//		_ = utils.PrintResponse[*model.Response[*model.ApplicationDto]](http.StatusOK, w, res)
//
// }
//
// func (a *applicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
//
//		user, err := GetUserFromContext(r.Context())
//		if err != nil {
//			panic(err)
//		}
//
//		apps := a.applicationRepository.GetAll(user.ID)
//
//		response := make([]*model.ApplicationDto, 0)
//		for _, app := range apps {
//			r := app.ToDTO()
//
//			response = append(response, r)
//		}
//
//		responseBody := model.NewResponse[[]*model.ApplicationDto](fmt.Sprintf("%d application(s) fetched successfully", len(apps)), response)
//
//		_ = utils.PrintResponse[*model.Response[[]*model.ApplicationDto]](http.StatusOK, w, responseBody)
//	}
//
//	func (a *applicationHandler) GenerateSecret(w http.ResponseWriter, r *http.Request) {
//		var request *model.CreateApplication
//
//		err := JsonToStruct(r.Body, &request)
//		if err != nil {
//			panic(errors.New("invalid request body"))
//		}
//
//		user, err := GetUserFromContext(r.Context())
//		if err != nil {
//			panic(errors.New("forbidden: not allowed to make this request"))
//		}
//
//		app, err := a.applicationRepository.GetByClientIdAndUserId(request.ClientId, user.ID)
//		if err != nil {
//			panic(err)
//		}
//
//		secret := uuid.NewString()
//		hashed, err := bcrypt.GenerateFromPassword([]byte(secret), 0)
//
//		if err != nil {
//			panic(errors.New("could not generate secret"))
//		}
//
//		appToUpdate := &model.Application{
//			ClientId:     app.ClientId,
//			ClientSecret: string(hashed),
//		}
//
//		_, err = a.applicationRepository.Update(appToUpdate)
//
//		if err != nil {
//			panic(errors.New("could not generate secret"))
//		}
//
//		applicationResponse := &model.ApplicationResponse{
//			ClientId:     request.ClientId,
//			ClientSecret: secret,
//			RedirectUrl:  app.RedirectUri,
//		}
//
//		response := model.NewResponse[*model.ApplicationResponse]("secret generated successfully", applicationResponse)
//
//		_ = utils.PrintResponse[*model.Response[*model.ApplicationResponse]](http.StatusOK, w, response)
//
// }
func NewApplicationHandler(applicationRepository kernel.Repository[model.Application], ctx kernel.Context) ApplicationHandler {
	return &applicationHandler{
		Context:    ctx,
		repository: applicationRepository,
	}
}

//
//func (a *applicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
//
//	var (
//		request *model.CreateApplication
//		user, _ = GetUserFromContext(r.Context())
//	)
//
//	err := JsonToStruct(r.Body, &request)
//	if err != nil {
//		panic(errors.New("invalid request body"))
//	}
//
//	if request.Name == "" {
//		panic(errors.New("name is required"))
//	}
//
//	alreadyExists, _ := a.applicationRepository.GetByNameAndUserId(request.Name, user.ID)
//	if alreadyExists != nil {
//		panic(errors.New("application with this name already exists"))
//	}
//
//	app := &model.Application{
//		Name:        request.Name,
//		RedirectUri: request.RedirectUri,
//		User:        *user,
//	}
//
//	err = a.applicationRepository.Create(app)
//	if err != nil {
//		panic(err)
//	}
//
//	response := model.NewResponse[*model.ApplicationDto]("application created successfully", app.ToDTO())
//
//	_ = utils.PrintResponse[*model.Response[*model.ApplicationDto]](http.StatusCreated, w, response)
//}
