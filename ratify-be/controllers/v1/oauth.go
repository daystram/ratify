package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
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
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// retrieve application
	var application models.Application
	if application, err = handlers.Handler.RetrieveApplication(authRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	flow := authRequest.Flow()
	switch flow {
	case constants.FlowAuthorizationCode, constants.FlowAuthorizationCodeWithPKCE:
		// verify user login
		var user models.User
		if user, err = handlers.Handler.AuthenticateUser(authRequest.UserLogin); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: "incorrect_credentials", Error: "incorrect username or password"})
			return
		}
		// verify request credentials
		// TODO: support comma-separated callback URLs
		if authRequest.RedirectURI != "" && authRequest.RedirectURI != application.CallbackURL {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "not allowed callback_uri"})
			return
		}
		// generate authorization code
		var authorizationCode string
		if authorizationCode, err = handlers.Handler.GenerateAuthorizationCode(authRequest, user.Subject); err != nil {
			c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed generating authorization_code"})
			return
		}
		// note code challenge
		if flow == constants.FlowAuthorizationCodeWithPKCE {
			if err = handlers.Handler.StoreCodeChallenge(authorizationCode, authRequest.PKCEAuthFields); err != nil {
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed storing code_challenge"})
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
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "unsupported authorization flow"})
		return
	}
}

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
			if tokenRequest.ClientSecret != application.ClientSecret {
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
	if err = bcrypt.CompareHashAndPassword([]byte(introspectRequest.ClientSecret), []byte(application.ClientSecret)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "invalid client_secret"})
		return
	}
	c.JSON(http.StatusOK, tokenInfo)
	return
}
