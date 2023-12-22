package handlers

import (
	_ "fmt"
)

//type TestHandler struct {
//	Repository kernel.Repository[model.Test]
//}

//func (h *TestHandler) Create(response http.ResponseWriter, request *http.Request) {
//
//	var req *model.Test
//
//	err := JsonToStruct(request.Body, &req)
//
//	if err != nil {
//		panic(errors.New("invalid request body"))
//	}
//
//	fmt.Println("test:", req)
//	err := h.Repository.Create(req)
//	if err != nil {
//		fmt.Println("error:" + err.Error())
//		panic(err)
//	}
//
//	_ = utils.PrintResponse[model.Test](201, response, *m)
//}
