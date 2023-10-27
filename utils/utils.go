package utils

import (
	"backend-server/model"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	bearerPrefix = "Bearer "
	UserKey      = "user"
)

func GenerateToken(client *model.Application, tokenSecret string) (string, error) {

	accessToken := model.AccessToken{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		Sub:       client.ClientId,
		Issuer:    fmt.Sprintf("http://%s:%s", getIp(), NewConfig().ServerPort()),
		Name:      client.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken)

	return token.SignedString([]byte(tokenSecret))
}

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

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	val := ctx.Value(UserKey)

	t, ok := val.(*model.User)

	if !ok {
		return nil, errors.New("could not verify user from context")
	}

	return t, nil
}
