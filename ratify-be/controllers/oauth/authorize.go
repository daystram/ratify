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
		var sessionID string
		var user models.User
		if authRequest.Username == "" && authRequest.Password == "" {
			// get session cookie
			sessionID, err = c.Cookie(constants.SessionIDCookieKey)
			if err != nil {
				c.SetCookie(constants.SessionIDCookieKey, "", -1, "/oauth", "", true, true)
				c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthIncorrectCredentials.Error(), Error: "invalid cookie"})
				return
			} else {
				// verify user session
				if user, sessionID, err = handlers.Handler.CheckSession(sessionID); err != nil {
					c.SetCookie(constants.SessionIDCookieKey, "", -1, "/oauth", "", true, true)
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: errors.ErrAuthIncorrectCredentials.Error(), Error: "invalid session_id"})
					return
				}
			}
		} else {
			// verify user login
			if user, sessionID, err = handlers.Handler.AuthenticateUser(authRequest.UserLogin); err != nil {
				if err == errors.ErrAuthIncorrectCredentials {
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: err.Error(), Error: "incorrect username or password"})
				} else if err == errors.ErrAuthEmailNotVerified {
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: err.Error(), Error: "email not verified"})
				} else {
					c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "failed logging in user"})
				}
				return
			}
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
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(constants.SessionIDCookieKey, sessionID, int(constants.SessionIDExpiry.Seconds()), "/oauth", "", true, true)
		c.JSON(http.StatusOK, gin.H{
			"data": fmt.Sprintf("%s?%s", strings.TrimSuffix(application.CallbackURL, "/"), param.Encode()),
		})
	default:
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "unsupported authorization flow"})
		return
	}
}

// @Summary Request logout
// @Tags oauth
// @Param user body datatransfers.LogoutRequest true "Logout request info"
// @Success 200 "OK"
// @Router /oauth/logout [POST]
func POSTLogout(c *gin.Context) {
	var err error
	// fetch request info
	var logoutRequest datatransfers.LogoutRequest
	if err = c.ShouldBindJSON(&logoutRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// retrieve application
	if _, err = handlers.Handler.RetrieveApplication(logoutRequest.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// revoke tokens
	if err = handlers.Handler.RevokeTokens(c.GetString(constants.UserSubjectKey), logoutRequest.ClientID, logoutRequest.Global); err != nil {
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
