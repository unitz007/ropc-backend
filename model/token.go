package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type claims map[string]interface{}

type AccessToken struct {
	claims claims
}

func NewAccessToken(subject string, expiry int, issuer string) *AccessToken {

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(time.Duration(expiry) * time.Minute).Unix()

	claims := claims{
		"iss": issuer,
		"sub": subject,
		"iat": issuedAt,
		"exp": expiresAt,
	}
	return &AccessToken{
		claims: claims,
	}
}

func (c *claims) Valid() error {
	return nil
}

func (t *AccessToken) Sign(secret string) string {
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &t.claims).SignedString([]byte(secret))
	if err != nil {
		panic(fmt.Errorf("failed to sign token: %v", err))
	}

	return signedToken

}

func (t *AccessToken) AddClaim(key, value string) *AccessToken {
	t.claims[key] = value
	return t
}
