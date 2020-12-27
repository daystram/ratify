package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) GenerateAuthorizationCode(application models.Application) (authorizationCode string, err error) {
	authorizationCode = utils.GenerateRandomString(constants.AuthorizationCodeLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode),
		application.ClientID, constants.AuthorizationCodeExpiry).Err(); err != nil {
		return "", errors.New(fmt.Sprintf("failed storing authorization code. %v", err))
	}
	return
}

func (m *module) ValidateAuthorizationCode(application models.Application, authorizationCode string) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)); result.Err() != nil {
		return errors.New(fmt.Sprintf("failed retrieving authorization code. %v", result.Err()))
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)) // immediate invalidation
	if result.Val() != application.ClientID {
		return errors.New("unregistered authorization code")
	}
	return
}

func (m *module) GenerateAccessRefreshToken(application models.Application) (accessToken, refreshToken string, err error) {
	accessToken = utils.GenerateRandomString(constants.AccessTokenLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyAccessToken, accessToken),
		application.ClientID, constants.AccessTokenExpiry).Err(); err != nil {
		return "", "", errors.New(fmt.Sprintf("failed storing access token. %v", err))
	}
	refreshToken = utils.GenerateRandomString(constants.RefreshTokenLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyRefreshToken, refreshToken),
		application.ClientID, constants.RefreshTokenExpiry).Err(); err != nil {
		return "", "", errors.New(fmt.Sprintf("failed storing refresh token. %v", err))
	}
	return
}

/*
Since an opaque token is used, applications must call an endpoint in ratify-be to validate the token on every request.
*/

func (m *module) StoreCodeChallenge(authorizationCode string, pkce datatransfers.PKCEAuthFields) (err error) {
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyCodeChallenge, authorizationCode),
		pkce.CodeChallengeMethod+constants.RDDelimiter+pkce.CodeChallenge,
		constants.AuthorizationCodeExpiry).Err(); err != nil {
		return errors.New(fmt.Sprintf("failed storing code challenge. %v", err))
	}
	return
}

func (m *module) ValidateCodeVerifier(authorizationCode string, pkce datatransfers.PKCETokenFields) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDKeyCodeChallenge, authorizationCode)); result.Err() != nil {
		return errors.New(fmt.Sprintf("failed retrieving code challenge. %v", err))
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDKeyCodeChallenge, authorizationCode))
	splitVal := strings.Split(result.Val(), constants.RDDelimiter)
	codeChallengeMethod, codeChallenge := splitVal[0], splitVal[1]
	switch codeChallengeMethod {
	case constants.PKCEChallengeMethodS256:
		hash := sha256.New()
		if _, err := hash.Write([]byte(pkce.CodeVerifier)); err != nil {
			return errors.New(fmt.Sprintf("failed hashing verifier. %v", err))
		}
		if codeChallenge != base64.RawURLEncoding.EncodeToString(hash.Sum([]byte{})) {
			return errors.New("failed verifying code challenge")
		}
	case constants.PKCEChallengeMethodPlain:
		if codeChallenge != pkce.CodeVerifier {
			return errors.New("failed verifying code challenge")
		}
	default:
		return errors.New(fmt.Sprintf("unsupported code challenge method %s", codeChallengeMethod))
	}
	return
}