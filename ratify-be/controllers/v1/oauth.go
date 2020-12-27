package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

// @Summary Request authorization
// @Tags oauth
// @Param user body datatransfers.AuthorizationRequest true "Authorization request info"
// @Success 200 "OK"
// @Router /oauth/authorize [POST]
func POSTAuthorize(c *gin.Context) {
	var err error
	// fetch request info
	var authRequest datatransfers.AuthorizationRequest
	if err = c.ShouldBindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}
	// retrieve application
	var application models.Application
	if application, err = handlers.Handler.RetrieveApplication(authRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "application not found"})
		return
	}
	flow := authRequest.Flow()
	switch flow {
	case constants.FlowAuthorizationCode, constants.FlowAuthorizationCodeWithPKCE:
		// verify user login
		if _, err = handlers.Handler.AuthenticateUser(authRequest.UserLogin); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "incorrect username or password"})
			return
		}
		// verify request credentials
		// TODO: support comma-separated callback URLs
		if authRequest.RedirectURI != "" && authRequest.RedirectURI != application.CallbackURL {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "not allowed callback URL"})
			return
		}
		// generate authorization code
		var authorizationCode string
		if authorizationCode, err = handlers.Handler.GenerateAuthorizationCode(application); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "failed generating authorization code"})
			return
		}
		// note code challenge
		if flow == constants.FlowAuthorizationCodeWithPKCE {
			if err = handlers.Handler.StoreCodeChallenge(authorizationCode, authRequest.PKCEAuthFields); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "failed storing code challenge"})
				return
			}
		}
		param, _ := query.Values(datatransfers.AuthorizationResponse{
			AuthorizationCode: authorizationCode,
			State:             authRequest.State,
		})
		c.JSON(http.StatusOK, gin.H{
			"data": fmt.Sprintf("%s?%s", strings.TrimSuffix(application.CallbackURL, "/"), param.Encode()),
		})
	default:
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: "unsupported authorization flow"})
		return
	}
}

// @Summary Request access (and refresh) tokens
// @Tags oauth
// @Param user body datatransfers.TokenRequest true "Token request info"
// @Success 200 "OK"
// @Router /oauth/token [POST]
func POSTToken(c *gin.Context) {
	var err error
	// fetch request info
	var tokenRequest datatransfers.TokenRequest
	if err = c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}
	// retrieve application
	var application models.Application
	if application, err = handlers.Handler.RetrieveApplication(tokenRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "application not found"})
		return
	}
	flow := tokenRequest.Flow()
	switch flow {
	case constants.FlowAuthorizationCode, constants.FlowAuthorizationCodeWithPKCE:
		// verify request credentials
		if err = handlers.Handler.ValidateAuthorizationCode(application, tokenRequest.Code); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "invalid authorization code"})
			return
		}
		if flow == constants.FlowAuthorizationCode {
			if tokenRequest.ClientSecret != application.ClientSecret {
				c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "invalid client secret"})
				return
			}
		}
		if flow == constants.FlowAuthorizationCodeWithPKCE {
			if err = handlers.Handler.ValidateCodeVerifier(tokenRequest.Code, tokenRequest.PKCETokenFields); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "failed verifying code challenge"})
				return
			}
		}
		// generate codes
		var accessToken, refreshToken string
		if accessToken, refreshToken, err = handlers.Handler.GenerateAccessRefreshToken(application); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "failed generating tokens"})
			return
		}
		c.JSON(http.StatusOK, datatransfers.Response{Data: datatransfers.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(constants.AccessTokenExpiry.Seconds()),
		}})
	default:
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: "unsupported grant_type"})
		return
	}
}