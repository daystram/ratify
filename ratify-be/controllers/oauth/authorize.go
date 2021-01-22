package oauth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/errors"
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
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// verify request credentials
	// TODO: support comma-separated callback URLs
	if authRequest.RedirectURI != "" && authRequest.RedirectURI != application.CallbackURL {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "not allowed callback_uri"})
		return
	}
	flow := authRequest.Flow()
	switch flow {
	case constants.FlowAuthorizationCode, constants.FlowAuthorizationCodeWithPKCE:
		var sessionID string
		var user models.User
		if authRequest.UseSession {
			// get session cookie
			sessionID, err = c.Cookie(constants.SessionIDCookieKey)
			if err != nil {
				c.SetCookie(constants.SessionIDCookieKey, "", -1, "/oauth", "", true, true)
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthIncorrectCredentials.Error(), Error: "invalid cookie"})
				return
			}
				// verify user session
				if user, sessionID, err = handlers.Handler.CheckSession(sessionID); err != nil {
					c.SetCookie(constants.SessionIDCookieKey, "", -1, "/oauth", "", true, true)
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthIncorrectCredentials.Error(), Error: "invalid session_id"})
					return
				}
		} else {
			// verify user login
			if user, sessionID, err = handlers.Handler.AuthenticateUser(authRequest.UserLogin); err != nil {
				if err == errors.ErrAuthIncorrectIdentifier {
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthIncorrectCredentials.Error(), Error: "incorrect credentials"})
				} else if err == errors.ErrAuthIncorrectCredentials {
					handlers.Handler.LogLogin(user, application, false, datatransfers.LogDetail{
						Scope:  constants.LogScopeOAuthAuthorize,
						Detail: utils.ParseUserAgent(c),
					})
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthIncorrectCredentials.Error(), Error: "incorrect credentials"})
				} else if err == errors.ErrAuthEmailNotVerified {
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthEmailNotVerified.Error(), Error: "email not verified"})
				} else if err == errors.ErrAuthMissingOTP { // proceed to MFA prompt
					c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Code: errors.ErrAuthMissingOTP.Error(), Error: "otp required"})
				} else {
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed logging in user"})
				}
				return
			}
		}
		// generate authorization code
		var authorizationCode string
		if authorizationCode, err = handlers.Handler.GenerateAuthorizationCode(authRequest, user.Subject, sessionID); err != nil {
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
		handlers.Handler.LogLogin(user, application, true, datatransfers.LogDetail{
			Scope:  constants.LogScopeOAuthAuthorize,
			Detail: utils.ParseUserAgent(c),
		})
		param, _ := query.Values(datatransfers.AuthorizationResponse{
			AuthorizationCode: authorizationCode,
			State:             authRequest.State,
		})
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(constants.SessionIDCookieKey, sessionID, int(constants.SessionIDExpiry.Seconds()), "/oauth", "", true, true)
		c.JSON(http.StatusOK, gin.H{
			"data": fmt.Sprintf("%s?%s", strings.TrimSuffix(application.CallbackURL, "/"), param.Encode()),
		})
		return
	default:
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "unsupported authorization flow"})
		return
	}
}

// @Summary Request logout
// @Tags oauth
// @Accept application/x-www-form-urlencoded
// @Param user body datatransfers.LogoutRequest true "Logout request info"
// @Success 200 "OK"
// @Router /oauth/logout [POST]
func POSTLogout(c *gin.Context) {
	var err error
	// fetch request info
	var logoutRequest datatransfers.LogoutRequest
	if err = c.ShouldBind(&logoutRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// introspect access_token
	var tokenInfo datatransfers.TokenIntrospection
	if tokenInfo, err = handlers.Handler.IntrospectAccessToken(logoutRequest.AccessToken); err != nil || !tokenInfo.Active {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: "invalid_token", Error: "invalid access_token"})
		return
	}
	// retrieve application
	if _, err = handlers.Handler.RetrieveApplication(logoutRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// revoke tokens
	if err = handlers.Handler.RevokeTokens(tokenInfo.Subject, logoutRequest.ClientID, logoutRequest.Global); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed revoking tokens"})
		return
	}
	if logoutRequest.Global {
		var sessionID string
		if sessionID, err = c.Cookie(constants.SessionIDCookieKey); err == nil {
			if err = handlers.Handler.ClearSession(sessionID); err != nil {
				log.Printf("failed clearing session. %v", err)
			}
		}
		c.SetCookie(constants.SessionIDCookieKey, "", -1, "/oauth", "", true, true)
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
