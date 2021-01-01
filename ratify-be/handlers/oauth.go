package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/daystram/ratify/ratify-be/config"
	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) GenerateAuthorizationCode(authRequest datatransfers.AuthorizationRequest, subject string) (authorizationCode string, err error) {
	authorizationCode = utils.GenerateRandomString(constants.AuthorizationCodeLength)
	key := fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)
	value := map[string]interface{}{
		"subject":   subject,
		"client_id": authRequest.ClientID,
		"scope":     authRequest.Scope,
	}
	if err = m.rd.HSet(context.Background(), key, value).Err(); err != nil {
		return "", errors.New(fmt.Sprintf("failed storing authorization_code. %v", err))
	}
	if err = m.rd.Expire(context.Background(), key, constants.AuthorizationCodeExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), key)
		return "", errors.New(fmt.Sprintf("failed setting authorization_code expiry. %v", err))
	}
	return
}

func (m *module) ValidateAuthorizationCode(application models.Application, authorizationCode string) (subject, scope string, err error) {
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)); result.Err() != nil {
		return "", "", errors.New(fmt.Sprintf("failed retrieving authorization_code. %v", result.Err()))
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)) // immediate invalidation
	var clientID string
	clientID, subject, scope = result.Val()["client_id"], result.Val()["subject"], result.Val()["scope"]
	if clientID != application.ClientID {
		return "", "", errors.New("unregistered authorization_code")
	}
	return
}

func (m *module) GenerateAccessRefreshToken(tokenRequest datatransfers.TokenRequest, subject string, withRefresh bool) (accessToken, refreshToken string, err error) {
	accessToken = utils.GenerateRandomString(constants.AccessTokenLength)
	accessTokenKey := fmt.Sprintf(constants.RDKeyAccessToken, accessToken)
	accessTokenValue := map[string]interface{}{
		"subject":   subject,
		"client_id": tokenRequest.ClientID,
	}
	if err = m.rd.HSet(context.Background(), accessTokenKey, accessTokenValue).Err(); err != nil {
		return "", "", errors.New(fmt.Sprintf("failed storing access_token. %v", err))
	}
	if err = m.rd.Expire(context.Background(), accessTokenKey, constants.AccessTokenExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), accessTokenKey)
		return "", "", errors.New(fmt.Sprintf("failed setting access_token expiry. %v", err))
	}
	if withRefresh {
		refreshTokenKey := fmt.Sprintf(constants.RDKeyAccessToken, accessToken)
		refreshTokenValue := map[string]interface{}{
			"subject":      subject,
			"client_id":    tokenRequest.ClientID,
			"access_token": accessToken,
		}
		refreshToken = utils.GenerateRandomString(constants.RefreshTokenLength)
		if err = m.rd.HSet(context.Background(), refreshTokenKey, refreshTokenValue).Err(); err != nil {
			return "", "", errors.New(fmt.Sprintf("failed storing refresh_token. %v", err))
		}
		if err = m.rd.Expire(context.Background(), refreshTokenKey, constants.RefreshTokenExpiry).Err(); err != nil {
			m.rd.Del(context.Background(), accessTokenKey)
			m.rd.Del(context.Background(), refreshTokenKey)
			return "", "", errors.New(fmt.Sprintf("failed setting refresh_token expiry. %v", err))
		}
	}
	return
}

/*
Since an opaque token is used, applications must call an endpoint in ratify-be to validate the token on every request.
*/

func (m *module) GenerateIDToken(clientID, subject string, scope []string) (idToken string, err error) {
	var user models.User
	if user, err = m.db.userOrmer.GetOneBySubject(subject); err != nil {
		return "", errors.New("cannot find user")
	}
	now := time.Now()
	expiry := time.Now().Add(constants.AuthenticationTimeout)
	claims := datatransfers.OpenIDClaims{
		Subject:   user.Subject,
		Issuer:    config.AppConfig.Domain,
		Audience:  clientID,
		ExpiresAt: expiry.Unix(),
		IssuedAt:  now.Unix(),
	}
	for _, s := range scope {
		switch s {
		case "profile":
			claims.Username = &user.Username
			claims.Superuser = &user.Superuser
			claims.GivenName = &user.GivenName
			claims.FamilyName = &user.FamilyName
			claims.UpdatedAt = &user.UpdatedAt
		case "email":
			claims.Email = &user.Email
			claims.EmailVerified = &user.EmailVerified
		}
	}
	return utils.GenerateJWT(claims)
}

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

func (m *module) IntrospectAccessToken(accessToken string) (tokenInfo datatransfers.TokenIntrospection, err error) {
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDKeyAccessToken, accessToken)); result.Err() != nil && err != redis.Nil {
		return datatransfers.TokenIntrospection{Active: false}, errors.New(fmt.Sprintf("failed retrieving access token. %v", result.Err()))
	}
	if err == redis.Nil {
		return datatransfers.TokenIntrospection{Active: false}, nil
	}
	return datatransfers.TokenIntrospection{
		Active:   true,
		ClientID: result.Val()["client_id"],
		Subject:  result.Val()["subject"],
		Scope:  result.Val()["scope"],
	}, nil
}
