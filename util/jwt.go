package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// use "secret" as the secret key
const SecretKey = "secret"

// GenerateJwt creates jwt and expires in 24 hours.
func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	// use "secret" as the secret key
	return claims.SignedString([]byte(SecretKey))
}

// ParseJwt parses cookie's jwt value and if valids, returns issuer which is user.id
func ParseJwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	// cast token claims into jwt.StandardClaims so that we can access to the claims.Issuer which is the user id
	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil
}
