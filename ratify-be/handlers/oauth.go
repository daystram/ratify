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

func (m *module) OAuthGenerateAuthorizationCode(authRequest datatransfers.AuthorizationRequest, subject, sessionID string) (authorizationCode string, err error) {
	authorizationCode = utils.GenerateRandomString(constants.AuthorizationCodeLength)
	key := fmt.Sprintf(constants.RDTemAuthorizationCode, authorizationCode)
	value := map[string]interface{}{
		"client_id":  authRequest.ClientID,
		"session_id": sessionID,
		"subject":    subject,
		"scope":      authRequest.Scope,
	}
	if err = m.rd.HSet(context.Background(), key, value).Err(); err != nil {
		return "", fmt.Errorf("failed storing authorization_code. %v", err)
	}
	if err = m.rd.Expire(context.Background(), key, constants.AuthorizationCodeExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), key)
		return "", fmt.Errorf("failed setting authorization_code expiry. %v", err)
	}
	return
}

func (m *module) OAuthValidateAuthorizationCode(application models.Application, authorizationCode string) (sessionID, subject, scope string, err error) {
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDTemAuthorizationCode, authorizationCode)); result.Err() != nil {
		return "", "", "", fmt.Errorf("failed retrieving authorization_code. %v", result.Err())
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemAuthorizationCode, authorizationCode)) // immediate invalidation
	var clientID string
	clientID, sessionID, subject, scope = result.Val()["client_id"], result.Val()["session_id"], result.Val()["subject"], result.Val()["scope"]
	if clientID != application.ClientID {
		return "", "", "", errors.New("unregistered authorization_code")
	}
	return
}

func (m *module) OAuthGenerateAccessToken(tokenRequest datatransfers.TokenRequest, sessionID, subject string, withRefresh bool) (accessToken, refreshToken string, err error) {
	accessToken = utils.GenerateRandomString(constants.AccessTokenLength)
	accessTokenKey := fmt.Sprintf(constants.RDTemAccessToken, accessToken)
	accessTokenValue := map[string]interface{}{
		"client_id":  tokenRequest.ClientID,
		"session_id": sessionID, //  would desync if session_id made rotating
		"subject":    subject,
	}
	if err = m.rd.HSet(context.Background(), accessTokenKey, accessTokenValue).Err(); err != nil {
		return "", "", fmt.Errorf("failed storing access_token. %v", err)
	}
	if err = m.rd.Expire(context.Background(), accessTokenKey, constants.AccessTokenExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), accessTokenKey)
		return "", "", fmt.Errorf("failed setting access_token expiry. %v", err)
	}
	if withRefresh {
		// NOTE: always check if embedded access_token has been revoked when attempting to renew using refresh_token
		refreshTokenKey := fmt.Sprintf(constants.RDTemAccessToken, accessToken)
		refreshTokenValue := map[string]interface{}{
			"client_id":    tokenRequest.ClientID,
			"subject":      subject,
			"access_token": accessToken,
		}
		refreshToken = utils.GenerateRandomString(constants.RefreshTokenLength)
		if err = m.rd.HSet(context.Background(), refreshTokenKey, refreshTokenValue).Err(); err != nil {
			return "", "", fmt.Errorf("failed storing refresh_token. %v", err)
		}
		if err = m.rd.Expire(context.Background(), refreshTokenKey, constants.RefreshTokenExpiry).Err(); err != nil {
			m.rd.Del(context.Background(), accessTokenKey)
			m.rd.Del(context.Background(), refreshTokenKey)
			return "", "", fmt.Errorf("failed setting refresh_token expiry. %v", err)
		}
	}
	return
}

/*
Since an opaque token is used, applications must call an endpoint in ratify-be to validate the token on every request.
*/

func (m *module) OAuthGenerateIDToken(clientID, subject string, scope []string) (idToken string, err error) {
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
			claims.CreatedAt = &user.CreatedAt
		case "email":
			claims.Email = &user.Email
			claims.EmailVerified = &user.EmailVerified
		}
	}
	return utils.GenerateJWT(claims)
}

func (m *module) OAuthStoreCodeChallenge(authorizationCode string, pkce datatransfers.PKCEAuthFields) (err error) {
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDTemCodeChallenge, authorizationCode),
		pkce.CodeChallengeMethod+constants.RDDelimiter+pkce.CodeChallenge,
		constants.AuthorizationCodeExpiry).Err(); err != nil {
		return fmt.Errorf("failed storing code challenge. %v", err)
	}
	return
}

func (m *module) OAuthValidateCodeVerifier(authorizationCode string, pkce datatransfers.PKCETokenFields) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDTemCodeChallenge, authorizationCode)); result.Err() != nil {
		return fmt.Errorf("failed retrieving code challenge. %v", err)
	}
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemCodeChallenge, authorizationCode))
	splitVal := strings.Split(result.Val(), constants.RDDelimiter)
	codeChallengeMethod, codeChallenge := splitVal[0], splitVal[1]
	switch codeChallengeMethod {
	case constants.PKCEChallengeMethodS256:
		hash := sha256.New()
		if _, err := hash.Write([]byte(pkce.CodeVerifier)); err != nil {
			return fmt.Errorf("failed hashing verifier. %v", err)
		}
		if codeChallenge != base64.RawURLEncoding.EncodeToString(hash.Sum([]byte{})) {
			return errors.New("failed verifying code challenge")
		}
	case constants.PKCEChallengeMethodPlain:
		if codeChallenge != pkce.CodeVerifier {
			return errors.New("failed verifying code challenge")
		}
	default:
		return fmt.Errorf("unsupported code challenge method %s", codeChallengeMethod)
	}
	return
}

func (m *module) OAuthIntrospectAccessToken(accessToken string) (tokenInfo datatransfers.TokenIntrospection, err error) {
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDTemAccessToken, accessToken)); result.Err() != nil {
		return datatransfers.TokenIntrospection{Active: false}, fmt.Errorf("failed retrieving access_token. %v", result.Err())
	}
	if result.Val()["client_id"] == "" {
		return datatransfers.TokenIntrospection{Active: false}, nil
	}
	return datatransfers.TokenIntrospection{
		Active:    true,
		ClientID:  result.Val()["client_id"],
		SessionID: result.Val()["session_id"],
		Subject:   result.Val()["subject"],
		Scope:     result.Val()["scope"],
	}, nil
}

func (m *module) OAuthRevokeAccessToken(accessToken string) (err error) {
	// NOTE: access_token listed in session_child will remain. Look at handlers.Handler.OAuthRevokeAllTokens
	return m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemAccessToken, accessToken)).Err()
}

// expensive flush
func (m *module) OAuthRevokeAllTokens(subject, clientID string, global bool) (err error) {
	var matches *redis.StringSliceCmd
	// revoke all access_token
	if matches = m.rd.Keys(context.Background(), constants.RDKeyAccessToken+"*"); matches.Err() != nil {
		return fmt.Errorf("failed retrieving access_token keys. %v", matches.Err())
	}
	for _, key := range matches.Val() {
		var result *redis.StringStringMapCmd
		if result = m.rd.HGetAll(context.Background(), key); result.Err() == nil {
			if result.Val()["subject"] == subject && (global || result.Val()["client_id"] == clientID) {
				_ = m.rd.Del(context.Background(), key)
			}
		}
	}
	// revoke all refresh_token
	if matches = m.rd.Keys(context.Background(), constants.RDKeyRefreshToken+"*"); matches.Err() != nil {
		return fmt.Errorf("failed retrieving refresh_token keys. %v", matches.Err())
	}
	for _, key := range matches.Val() {
		var result *redis.StringStringMapCmd
		if result = m.rd.HGetAll(context.Background(), key); result.Err() == nil {
			if result.Val()["subject"] == subject && (global || result.Val()["client_id"] == clientID) {
				_ = m.rd.Del(context.Background(), key)
			}
		}
	}
	return
}
