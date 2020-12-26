package utils

import (
	"crypto/rand"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/daystram/ratify/ratify-be/config"
	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
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

func GenerateJWT(user models.User) (string, error) {
	now := time.Now()
	expiry := time.Now().Add(constants.AuthenticationTimeout)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, datatransfers.JWTClaims{
		Subject:     user.Subject,
		IsSuperuser: user.Superuser,
		ExpiresAt:   expiry.Unix(),
		IssuedAt:    now.Unix(),
	})
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
