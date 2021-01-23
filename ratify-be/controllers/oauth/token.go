package oauth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

// @Summary Request access (and refresh) tokens
// @Tags oauth
// @Accept application/x-www-form-urlencoded
// @Param user body datatransfers.TokenRequest true "Token request info"
// @Success 200 "OK"
// @Router /oauth/token [POST]
func POSTToken(c *gin.Context) {
	var err error
	// fetch request info
	var tokenRequest datatransfers.TokenRequest
	if err = c.ShouldBind(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// retrieve application
	var application models.Application
	if application, err = handlers.Handler.ApplicationGetOneByClientID(tokenRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	flow := tokenRequest.Flow()
	switch flow {
	case constants.FlowAuthorizationCode, constants.FlowAuthorizationCodeWithPKCE:
		// verify request credentials
		var sessionID, subject, scope string
		if sessionID, subject, scope, err = handlers.Handler.OAuthValidateAuthorizationCode(application, tokenRequest.Code); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "invalid authorization_code"})
			return
		}
		if flow == constants.FlowAuthorizationCode {
			if err = bcrypt.CompareHashAndPassword([]byte(application.ClientSecret), []byte(tokenRequest.ClientSecret)); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "invalid client_secret"})
				return
			}
		}
		if flow == constants.FlowAuthorizationCodeWithPKCE {
			if err = handlers.Handler.OAuthValidateCodeVerifier(tokenRequest.Code, tokenRequest.PKCETokenFields); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed verifying code_challenge"})
				return
			}
		}
		// generate codes
		var accessToken, refreshToken, idToken string
		if utils.HasOpenIDScope(scope) {
			if idToken, err = handlers.Handler.OAuthGenerateIDToken(application.ClientID, subject, strings.Split(scope, " ")); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed generating tokens"})
				return
			}
		}
		if accessToken, refreshToken, err = handlers.Handler.OAuthGenerateAccessToken(tokenRequest, sessionID, subject, flow == constants.FlowAuthorizationCode); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed generating tokens"})
			return
		}
		// spawn child session
		if err = handlers.Handler.SessionAddChild(sessionID, accessToken); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed spawning child session"})
			return
		}
		handlers.Handler.LogInsertAuthorize(application, true, datatransfers.LogDetail{})
		c.JSON(http.StatusOK, datatransfers.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			IDToken:      idToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(constants.AccessTokenExpiry.Seconds()),
		})
		return
	default:
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "unsupported grant_type"})
		return
	}
}

// @Summary Introspect token
// @Tags oauth
// @Accept application/x-www-form-urlencoded
// @Param user body datatransfers.IntrospectRequest true "Token request info"
// @Success 200 "OK"
// @Router /oauth/introspect [POST]
func POSTIntrospect(c *gin.Context) {
	var err error
	// fetch request info
	var introspectRequest datatransfers.IntrospectRequest
	if err = c.ShouldBind(&introspectRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// verify client_id and client_secret
	var application models.Application
	if application, err = handlers.Handler.ApplicationGetOneByClientID(introspectRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(application.ClientSecret), []byte(introspectRequest.ClientSecret)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "invalid client_secret"})
		return
	}
	// introspect
	// TODO: allow introspecting other token types
	var tokenInfo datatransfers.TokenIntrospection
	if tokenInfo, err = handlers.Handler.OAuthIntrospectAccessToken(introspectRequest.Token); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed introspecting token"})
		return
	}
	c.JSON(http.StatusOK, tokenInfo)
	return
}

// @Summary Get user info from access_token
// @Tags oauth
// @Security BearerAuth
// @Success 200 "OK"
// @Router /oauth/userinfo [GET]
func GETUserInfo(c *gin.Context) {
	var err error
	var user models.User
	if user, err = handlers.Handler.UserGetOneBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.UserInfo{
		FamilyName:    user.FamilyName,
		GivenName:     user.GivenName,
		Subject:       user.Subject,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	})
	return
}
