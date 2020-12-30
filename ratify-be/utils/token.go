package utils

import (
	"crypto/rand"

	"github.com/dgrijalva/jwt-go"

	"github.com/daystram/ratify/ratify-be/config"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	byteString := make([]byte, length)
	_, _ = rand.Read(byteString)
	for i, b := range byteString {
		byteString[i] = letterBytes[b%byte(len(letterBytes))]
	}
	return string(byteString)
}

func GenerateJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
