package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"ropc-backend/model"
	"testing"

	"gorm.io/gorm"
)

const (
	jsonDecodeError = "Couldn't decode JSON: "
	UserKey         = "user"
)

func JsonToStruct[T any](r io.ReadCloser, t T) error {
	err := json.NewDecoder(r).Decode(t)
	if err != nil {
		log.Println("error = ", err)
		return errors.New(jsonDecodeError + err.Error())
	}

	return nil
}

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	val := ctx.Value(UserKey)

	t, ok := val.(*model.User)

	if !ok {
		return nil, errors.New("could not verify user from context")
	}

	return t, nil
}

func BuildTestRequest(t testing.TB, body io.Reader) (req *httptest.ResponseRecorder, res *http.Request) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPut, "http://localhost:0909/apps", body)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request = request.WithContext(context.WithValue(request.Context(), UserKey, &model.User{Model: gorm.Model{ID: uint(2)}}))
	response := httptest.NewRecorder()

	return response, request
}
