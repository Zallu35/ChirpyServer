package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIN time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-acces",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIN)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.Nil, err
	}
	userIDstr, er := token.Claims.GetSubject()
	if er != nil {
		return uuid.Nil, er
	}
	userID, e := uuid.Parse(userIDstr)
	if e != nil {
		return uuid.Nil, e
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return bearer, errors.New("no Authorization header found")
	}
	token := strings.TrimPrefix(bearer, "Bearer ")
	return token, nil
}
