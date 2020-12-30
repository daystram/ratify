package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

// @Summary Get user detail
// @Tags user
// @Security BearerAuth
// @Param user body datatransfers.UserSignup true "User signup info"
// @Success 200 "OK"
// @Router /api/v1/user [GET]
func GETUser(c *gin.Context) {
	var err error
	var user models.User
	if user, err = handlers.Handler.RetrieveUserByUsername(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: datatransfers.UserInfo{
		FamilyName: user.FamilyName,
		GivenName:  user.GivenName,
		Subject:    user.Subject,
		Username:   user.Username,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
	}})
	return
}

// @Summary Register user
// @Tags user
// @Param user body datatransfers.UserSignup true "User signup info"
// @Success 200 "OK"
// @Router /api/v1/user [POST]
func POSTRegister(c *gin.Context) {
	var err error
	var user datatransfers.UserSignup
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	if user, _ := handlers.Handler.RetrieveUserByUsername(user.Username); user.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: "username_exists", Error: "username already used"})
		return
	}
	if user, _ := handlers.Handler.RetrieveUserByEmail(user.Email); user.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: "email_exists", Error: "email already used"})
		return
	}
	if err = handlers.Handler.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed registering user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Update user
// @Tags user
// @Security BearerAuth
// @Param user body datatransfers.UserSignup true "User update info"
// @Success 200 "OK"
// @Router /api/v1/user [PUT]
func PUTUser(c *gin.Context) {
	var err error
	var user datatransfers.UserUpdate
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	if user, _ := handlers.Handler.RetrieveUserByEmail(user.Email); user.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: "email_exists", Error: "email already used"})
		return
	}
	if err = handlers.Handler.UpdateUser(c.GetString(constants.IsAuthenticatedKey), user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
