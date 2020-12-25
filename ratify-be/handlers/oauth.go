package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) GenerateAuthorizationCode(application models.Application) (authorizationCode string, err error) {
	authorizationCode = utils.GenerateRandomString(constants.AuthorizationCodeLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.AuthorizationCodeKey, authorizationCode),
		application.ClientID, constants.AuthorizationCodeExpiry).Err(); err != nil {
		return "", errors.New(fmt.Sprintf("failed storing authorization code. %v", err))
	}
	return
}

func (m *module) ValidateAuthorizationCode(application models.Application, authorizationCode string) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.AuthorizationCodeKey, authorizationCode)); result.Err() != nil {
		return errors.New(fmt.Sprintf("failed retrieving authorization code. %v", result.Err()))
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.AuthorizationCodeKey, authorizationCode)) // immediate invalidation
	if result.Val() != application.ClientID {
		return errors.New("unregistered authorization code")
	}
	return
}

func (m *module) GenerateAccessRefreshToken(application models.Application) (accessToken, refreshToken string, err error) {
	accessToken = utils.GenerateRandomString(constants.AccessTokenLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.AccessTokenKey, accessToken),
		application.ClientID, constants.AccessTokenExpiry).Err(); err != nil {
		return "", "", errors.New(fmt.Sprintf("failed storing access token. %v", err))
	}
	refreshToken = utils.GenerateRandomString(constants.RefreshTokenLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RefreshTokenKey, refreshToken),
		application.ClientID, constants.RefreshTokenExpiry).Err(); err != nil {
		return "", "", errors.New(fmt.Sprintf("failed storing refresh token. %v", err))
	}
	return
}
/*
Since an opaque token is used, applications must call an endpoint in ratify-be to validate the token on every request.
 */
