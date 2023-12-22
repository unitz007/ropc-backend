package kernel

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"ropc-backend/model"
	"testing"
)

type contextMock struct{}
type loggerMock struct{}

func (l loggerMock) Info(v string)  {}
func (l loggerMock) Warn(v string)  {}
func (l loggerMock) Error(v string) {}
func (l loggerMock) Fatal(v string) {}

func (m contextMock) Database() Database {
	return nil
}

func (m contextMock) Logger() Logger {
	return loggerMock{}
}

func (m contextMock) Router() Router {
	return nil
}
func Test_PanicHandler(t *testing.T) {

	scenarios := []struct {
		name         string
		handlerFunc  func(w http.ResponseWriter, r *http.Request)
		expectedCode int
		expectedMsg  string
	}{
		{
			"should be 400 status code",
			func(w http.ResponseWriter, r *http.Request) {
				panic(errors.New("dummy message"))
			},
			http.StatusBadRequest,
			"dummy message",
		},
		{
			"should be 500 status code",
			func(w http.ResponseWriter, r *http.Request) {
				panic("test dummy")
			},
			http.StatusInternalServerError,
			"Internal Server Error",
		},
		{
			"should be 404 status code",
			func(w http.ResponseWriter, r *http.Request) {
				panic(NewError(http.StatusNotFound, "dummy message"))
			},
			http.StatusNotFound,
			"dummy message",
		},
	}

	for _, test := range scenarios {
		t.Run(test.name, func(t *testing.T) {
			middlewareTest(t, test.handlerFunc, test.expectedCode, test.expectedMsg)
		})
	}
}

func middlewareTest(t testing.TB, handlerFunc func(http.ResponseWriter, *http.Request), expectedCode int, expectedMsg string) *httptest.ResponseRecorder {
	t.Helper()
	panicHandler := NewMiddleware(contextMock{}).PanicHandler(handlerFunc)
	response := httptest.NewRecorder()

	panicHandler(response, nil)

	if response.Code != expectedCode {
		t.Errorf("got %d, want %d", response.Code, expectedCode)
	}

	var responseBody model.Response[any]

	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal("could  not unmarshal")
	}

	if responseBody.Message != expectedMsg {
		t.Errorf("expected \"%s\" but got \"%s\"", expectedMsg, responseBody.Message)
	}

	return response
}
