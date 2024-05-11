package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type AuthorizationHandler interface {
	IDK(w http.ResponseWriter, r *http.Request)
}

type authorizationHandler struct{}

func NewAuthorizationHandler() AuthorizationHandler {
	return &authorizationHandler{}
}

func (h authorizationHandler) IDK(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	//fmt.Println(code)
	t, _ := jwt.Parse(code, nil)

	s := t.Claims.(jwt.MapClaims)["sub"].(string)
	fmt.Println(s)

	err := os.Setenv("IDK_USER_ID", s)
	if err != nil {
		return
	}

	html := `<html>
			<head><title>CLI Authorization</title></head>
			<body>
				<h1>Authorization Successful</h1>
			</body>
			</html>`

	_, _ = w.Write([]byte(html))
}
