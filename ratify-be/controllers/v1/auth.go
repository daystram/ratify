package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
)

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
