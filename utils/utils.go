package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"ropc-backend/model"
	"strings"

	"github.com/golang-jwt/jwt"
)

const (
	bearerPrefix = "Bearer "
)

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

func PrintResponseNew[T any](responseWriter http.ResponseWriter, statusCode int, message string, payload T) error {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	response := model.NewResponse[T](message, payload)

	err := json.NewEncoder(responseWriter).Encode(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		_, err = responseWriter.Write([]byte("Invalid response"))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetIssuerUri(conf Config) string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return "http://" + ip + ":" + conf.ServerPort()
}
