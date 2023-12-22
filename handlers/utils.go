package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"ropc-backend/utils"
	"testing"

	"gorm.io/gorm"
)

const (
	jsonDecodeError = "Couldn't decode JSON: "
)

func JsonToStruct[T any](r io.ReadCloser, t T) error {
	err := json.NewDecoder(r).Decode(t)
	if err != nil {
		log.Println("error = ", err)
		return errors.New(jsonDecodeError + err.Error())
	}

	return nil
}

func GetUserFromContext(ctx context.Context) *model.User {
	val := ctx.Value(utils.UserKey)

	t, ok := val.(*model.User)

	if !ok {
		panic(kernel.NewError(http.StatusForbidden, "could not verify user from context"))
	}

	return t
}

func BuildTestRequest(t testing.TB, body io.Reader) (req *httptest.ResponseRecorder, res *http.Request) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps", body)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request = request.WithContext(context.WithValue(request.Context(), utils.UserKey, &model.User{Model: gorm.Model{ID: uint(2)}}))
	response := httptest.NewRecorder()

	return response, request
}
