package handlers

import (
	"errors"
	_ "fmt"
	"net/http"
	"ropc-backend/kernel"
	"ropc-backend/model"
)

type TestHandler struct {
	Repository kernel.Repository[model.Test]
}

func (h *TestHandler) Create(response http.ResponseWriter, request *http.Request) {

	var req model.Test

	err := JsonToStruct(request.Body, &req)

	if err != nil {
		panic(errors.New("invalid request body"))
	}

	http.Error(response, "ndkjnekj", 200)

	return
	//fmt.Println("test:", req)
	//m, err := h.Repository.Create(req)
	//if err != nil {
	//	fmt.Println("error:" + err.Error())
	//	panic(err)
	//}err

	//_ = utils.PrintResponse[model.Test](201, response, *m)
}
