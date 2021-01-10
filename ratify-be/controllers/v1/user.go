package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/errors"
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
	if user, err = handlers.Handler.RetrieveUserBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: datatransfers.UserInfo{
		FamilyName:    user.FamilyName,
		GivenName:     user.GivenName,
		Subject:       user.Subject,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
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
	// fetch signup info
	var signup datatransfers.UserSignup
	if err = c.ShouldBindJSON(&signup); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// check email and username
	if checkUser, _ := handlers.Handler.RetrieveUserByUsername(signup.Username); checkUser.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: errors.ErrUserUsernameExists.Error(), Error: "username already used"})
		return
	}
	if checkUser, _ := handlers.Handler.RetrieveUserByEmail(signup.Email); checkUser.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: errors.ErrUserEmailExists.Error(), Error: "email already used"})
		return
	}
	// register user
	var user models.User
	if user.Subject, err = handlers.Handler.RegisterUser(signup); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed registering user"})
		return
	}
	// send verification email
	if user, err = handlers.Handler.RetrieveUserBySubject(user.Subject); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed retrieving new user"})
		return
	}
	if err = handlers.Handler.SendVerificationEmail(user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed sending verification email"})
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
	if user, _ := handlers.Handler.RetrieveUserByEmail(user.Email); user.Subject != "" && user.Subject != c.GetString(constants.UserSubjectKey) {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: errors.ErrUserEmailExists.Error(), Error: "email already used"})
		return
	}
	if err = handlers.Handler.UpdateUser(c.GetString(constants.UserSubjectKey), user); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Verify user email
// @Tags user
// @Param user body datatransfers.UserVerify true "User verification info"
// @Success 200 "OK"
// @Router /api/v1/user/verify [POST]
func POSTVerify(c *gin.Context) {
	var err error
	// fetch verification info
	var verify datatransfers.UserVerify
	if err = c.ShouldBindJSON(&verify); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// verify user
	if err := handlers.Handler.VerifyUser(verify.Token); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "failed verifying user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
