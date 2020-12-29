package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
)

// @Summary Login user
// @Tags auth
// @Param user body datatransfers.UserLogin true "User login info"
// @Success 200 "OK"
// @Router /api/v1/auth/login [POST]
func POSTLogin(c *gin.Context) {
	var err error
	var user datatransfers.UserLogin
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}
	var token string
	if token, err = handlers.Handler.AuthenticateUser(user); err != nil {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "incorrect username or password"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.Response{Data: fmt.Sprintf("Bearer %s", token)})
	return
}

// @Summary Register user
// @Tags auth
// @Param user body datatransfers.UserSignup true "User signup info"
// @Success 200 "OK"
// @Router /api/v1/auth/register [POST]
func POSTRegister(c *gin.Context) {
	var err error
	var user datatransfers.UserSignup
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	if err = handlers.Handler.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed registering user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
