package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

const (
	bearerPrefix = "Bearer "
)

//func GenerateToken[T model.TokenProperties](t T, tokenSecret string) (string, error) {
//
//	accessToken := model.AccessToken{
//		IssuedAt:  time.Now().Unix(),
//		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
//		Sub:       t.Subject(),
//		Issuer:    fmt.Sprintf("http://%s:%s", getIp(), NewConfig().ServerPort()),
//		Name:      client.Name,
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken)
//
//	return token.SignedString([]byte(tokenSecret))
//}

func ValidateToken(token, tokenSecret string) (jwt.MapClaims, error) {
	token = strings.TrimPrefix(token, bearerPrefix)
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, err
}

func PrintResponse[T any](statusCode int, res http.ResponseWriter, payload T) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	err := json.NewEncoder(res).Encode(payload)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, err = res.Write([]byte("Invalid response"))
		if err != nil {
			return err
		}
	}

	return nil
}
