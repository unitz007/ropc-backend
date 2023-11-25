package model

import (
	"fmt"
	"net"
	"ropc-backend/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

type claims map[string]interface{}

type AccessToken struct {
	claims claims
	config utils.Config
}

func NewAccessToken(subject string, config utils.Config) *AccessToken {

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(time.Duration(config.TokenExpiry()) * time.Minute).Unix()

	claims := claims{
		"iss": fmt.Sprintf("http://%s:%s", getIp(), config.ServerPort()),
		"sub": subject,
		"iat": issuedAt,
		"exp": expiresAt,
	}
	return &AccessToken{
		claims: claims,
		config: config,
	}
}

func (c *claims) Valid() error {
	return nil
}

func (t *AccessToken) Sign() string {
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &t.claims).SignedString([]byte(t.config.TokenSecret()))
	if err != nil {
		panic(fmt.Errorf("failed to sign token: %v", err))
	}

	return signedToken

}

func (t *AccessToken) AddClaim(key, value string) *AccessToken {
	t.claims[key] = value
	return t
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
