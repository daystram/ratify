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

func (m *module) GenerateAuthorizationCode(application models.Application, subject string) (authorizationCode string, err error) {
	authorizationCode = utils.GenerateRandomString(constants.AuthorizationCodeLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode),
		application.ClientID+constants.RDDelimiter+subject, constants.AuthorizationCodeExpiry).Err(); err != nil {
		return "", errors.New(fmt.Sprintf("failed storing authorization code. %v", err))
	}
	return
}

func (m *module) ValidateAuthorizationCode(application models.Application, authorizationCode string) (subject string, err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)); result.Err() != nil {
		return "", errors.New(fmt.Sprintf("failed retrieving authorization code. %v", result.Err()))
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDKeyAuthorizationCode, authorizationCode)) // immediate invalidation
	splitVal := strings.Split(result.Val(), constants.RDDelimiter)
	if splitVal[0] != application.ClientID {
		return "", errors.New("unregistered authorization code")
	}
	return splitVal[1], nil
}

func (m *module) GenerateAccessRefreshToken(application models.Application, subject string, withRefresh bool) (accessToken, refreshToken string, err error) {
	accessToken = utils.GenerateRandomString(constants.AccessTokenLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyAccessToken, accessToken),
		application.ClientID+constants.RDDelimiter+subject, constants.AccessTokenExpiry).Err(); err != nil {
		return "", "", errors.New(fmt.Sprintf("failed storing access token. %v", err))
	}
	if withRefresh {
		refreshToken = utils.GenerateRandomString(constants.RefreshTokenLength)
		if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDKeyRefreshToken, refreshToken),
			application.ClientID+constants.RDDelimiter+subject, constants.RefreshTokenExpiry).Err(); err != nil {
			return "", "", errors.New(fmt.Sprintf("failed storing refresh token. %v", err))
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
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDKeyAccessToken, accessToken)); result.Err() != nil && err != redis.Nil {
		return datatransfers.TokenIntrospection{}, errors.New(fmt.Sprintf("failed retrieving access token. %v", result.Err()))
	}
	if err == redis.Nil {
		return datatransfers.TokenIntrospection{Active: false}, nil
	}
	splitVal := strings.Split(result.Val(), constants.RDDelimiter)
	return datatransfers.TokenIntrospection{
		Active:   true,
		ClientID: splitVal[0],
		Subject:  splitVal[1],
	}, nil
}
