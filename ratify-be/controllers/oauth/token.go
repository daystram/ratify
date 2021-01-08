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
	if application, err = handlers.Handler.RetrieveApplication(tokenRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	flow := tokenRequest.Flow()
	switch flow {
	case constants.FlowAuthorizationCode, constants.FlowAuthorizationCodeWithPKCE:
		// verify request credentials
		var subject, scope string
		if subject, scope, err = handlers.Handler.ValidateAuthorizationCode(application, tokenRequest.Code); err != nil {
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
			if err = handlers.Handler.ValidateCodeVerifier(tokenRequest.Code, tokenRequest.PKCETokenFields); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed verifying code_challenge"})
				return
			}
		}
		// generate codes
		var accessToken, refreshToken, idToken string
		if utils.HasOpenIDScope(scope) {
			if idToken, err = handlers.Handler.GenerateIDToken(application.ClientID, subject, strings.Split(scope, " ")); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed generating tokens"})
				return
			}
		}
		if accessToken, refreshToken, err = handlers.Handler.GenerateAccessRefreshToken(tokenRequest, subject, flow == constants.FlowAuthorizationCode); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed generating tokens"})
			return
		}
		c.JSON(http.StatusOK, datatransfers.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			IDToken:      idToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(constants.AccessTokenExpiry.Seconds()),
		})
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
	// retrieve details
	// TODO: allow introspecting other token types
	var tokenInfo datatransfers.TokenIntrospection
	if tokenInfo, err = handlers.Handler.IntrospectAccessToken(introspectRequest.Token); err != nil || !tokenInfo.Active {
		c.JSON(http.StatusOK, tokenInfo)
		return
	}
	// verify client_secret
	var application models.Application
	if application, err = handlers.Handler.RetrieveApplication(tokenInfo.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(application.ClientSecret), []byte(introspectRequest.ClientSecret)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "invalid client_secret"})
		return
	}
	c.JSON(http.StatusOK, tokenInfo)
	return
}
