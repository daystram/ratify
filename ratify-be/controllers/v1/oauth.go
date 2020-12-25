package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

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
	switch authRequest.ResponseType {
	case constants.ResponseTypeCode:
		// verify user login
		if _, err = handlers.Handler.AuthenticateUser(authRequest.UserLogin); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "incorrect username or password"})
			return
		}
		// verify request credentials
		// TODO: support comma-separated callback URLs
		if authRequest.RedirectURI != application.CallbackURL {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "not allowed callback URL"})
			return
		}
		// generate authorization code
		var authorizationCode string
		if authorizationCode, err = handlers.Handler.GenerateAuthorizationCode(application); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "failed generating authorization code"})
			return
		}
		c.JSON(http.StatusOK, datatransfers.Response{Data: datatransfers.AuthorizationResponse{
			AuthorizationCode: authorizationCode,
			State:             authRequest.State,
		}})
	default:
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: "unsupported response_type"})
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
	switch tokenRequest.GrantType {
	case constants.GrantTypeAuthorizationCode:
		// verify request credentials
		if tokenRequest.ClientSecret != application.ClientSecret {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "invalid client secret"})
			return
		}
		if err = handlers.Handler.ValidateAuthorizationCode(application, tokenRequest.Code); err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "invalid authorization code"})
			return
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
