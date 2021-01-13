package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/errors"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

// @Summary Enable TOTP
// @Tags mfa
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/mfa/enable [POST]
func POSTEnableTOTP(c *gin.Context) {
	var err error
	// get user
	var user models.User
	if user, err = handlers.Handler.RetrieveUserBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed retrieving user"})
		return
	}
	var uri string
	if uri, err = handlers.Handler.EnableTOTP(user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: fmt.Sprintf("failed enabling totp. %v", err)})
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: uri})
	return
}

// @Summary Confirm TOTP
// @Tags mfa
// @Security BearerAuth
// @Param user body datatransfers.TOTPRequest true "User otp info"
// @Success 200 "OK"
// @Router /api/v1/mfa/confirm [POST]
func POSTConfirmTOTP(c *gin.Context) {
	var err error
	// fetch otp
	var totp datatransfers.TOTPRequest
	if err = c.ShouldBindJSON(&totp); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// get user
	var user models.User
	if user, err = handlers.Handler.RetrieveUserBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed retrieving user"})
		return
	}
	// confirm TOTP
	if err = handlers.Handler.ConfirmTOTP(totp.OTP, user); err != nil {
		if err == errors.ErrAuthIncorrectCredentials {
			c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Code: err.Error(), Error: "incorrect otp"})
		} else {
			c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: fmt.Sprintf("failed confirming totp. %v", err)})
		}
		return
	}
	handlers.Handler.LogUser(user, true, datatransfers.LogDetail{
		Scope:  constants.LogScopeUserMFA,
		Detail: true,
	})
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Disable TOTP
// @Tags mfa
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/mfa/disable [POST]
func POSTDisableTOTP(c *gin.Context) {
	var err error
	// get user
	var user models.User
	if user, err = handlers.Handler.RetrieveUserBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed retrieving user"})
		return
	}
	// disable TOTP
	if err = handlers.Handler.DisableTOTP(user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: fmt.Sprintf("failed disabling totp. %v", err)})
		return
	}
	handlers.Handler.LogUser(user, true, datatransfers.LogDetail{
		Scope:  constants.LogScopeUserMFA,
		Detail: false,
	})
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
