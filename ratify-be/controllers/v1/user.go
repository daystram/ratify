package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

func GETUser(c *gin.Context) {
	var err error
	var userInfo datatransfers.UserInfo
	if err = c.ShouldBindUri(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	var user models.User
	if user, err = handlers.Handler.RetrieveUserByUsername(userInfo.Username); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: datatransfers.UserInfo{
		Subject:   user.Subject,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}})
	return
}

func PUTUser(c *gin.Context) {
	var err error
	var user datatransfers.UserUpdate
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	if err = handlers.Handler.UpdateUser(c.GetString(constants.IsAuthenticatedKey), user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
