package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/errors"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

// @Summary Get user detail
// @Tags user
// @Security BearerAuth
// @Param user body datatransfers.UserSignup true "User signup info"
// @Success 200 "OK"
// @Router /api/v1/user/{subject} [GET]
func GETUserDetail(c *gin.Context) {
	var err error
	// check superuser
	var user models.User
	user.Subject = strings.TrimPrefix(c.Param("subject"), "/")
	if !c.GetBool(constants.IsSuperuserKey) && (user.Subject != c.GetString(constants.UserSubjectKey)) {
		c.JSON(http.StatusUnauthorized, datatransfers.APIResponse{Error: "access unauthorized"})
		return
	}
	// retrieve user
	if user, err = handlers.Handler.UserGetOneBySubject(user.Subject); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: datatransfers.UserInfo{
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		Subject:       user.Subject,
		Superuser:     user.Superuser,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		MFAEnabled:    user.EnabledTOTP(),
		CreatedAt:     user.CreatedAt,
		SignInCount:   user.LoginCount,
		LastSignIn:    user.LastLogin,
	}})
	return
}

// @Summary Get all users
// @Tags user
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/user [GET]
func GETUserList(c *gin.Context) {
	var err error
	// get all users
	var users []models.User
	if users, err = handlers.Handler.UserGetAll(); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "cannot retrieve users"})
		return
	}
	var usersResponse []datatransfers.UserInfo
	for _, user := range users {
		usersResponse = append(usersResponse, datatransfers.UserInfo{
			Subject:       user.Subject,
			GivenName:     user.GivenName,
			FamilyName:    user.FamilyName,
			Username:      user.Username,
			EmailVerified: user.EmailVerified,
			MFAEnabled:    user.EnabledTOTP(),
			CreatedAt:     user.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: usersResponse})
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
	if checkUser, _ := handlers.Handler.UserGetOneByUsername(signup.Username); checkUser.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: errors.ErrUserUsernameExists.Error(), Error: "username already used"})
		return
	}
	if checkUser, _ := handlers.Handler.UserGetOneByEmail(signup.Email); checkUser.Subject != "" {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: errors.ErrUserEmailExists.Error(), Error: "email already used"})
		return
	}
	// register user
	var user models.User
	if user.Subject, err = handlers.Handler.AuthRegister(signup); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed registering user"})
		return
	}
	// send verification email
	if user, err = handlers.Handler.UserGetOneBySubject(user.Subject); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed retrieving new user"})
		return
	}
	if err = handlers.Handler.MailerSendEmailVerification(user); err != nil {
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
	var update datatransfers.UserUpdate
	if err = c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	var user models.User
	if user, _ = handlers.Handler.UserGetOneByEmail(update.Email); user.Subject != "" && user.Subject != c.GetString(constants.UserSubjectKey) {
		c.JSON(http.StatusConflict, datatransfers.APIResponse{Code: errors.ErrUserEmailExists.Error(), Error: "email already used"})
		return
	}
	if err = handlers.Handler.UserUpdate(c.GetString(constants.UserSubjectKey), update); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating user"})
		return
	}
	handlers.Handler.LogInsertUser(user, true, datatransfers.LogDetail{Scope: constants.LogScopeUserProfile})
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Update user password
// @Tags user
// @Param user body datatransfers.UserUpdatePassword true "User update password info"
// @Success 200 "OK"
// @Router /api/v1/user/password [PUT]
func PUTUserPassword(c *gin.Context) {
	var err error
	// fetch verification info
	var password datatransfers.UserUpdatePassword
	if err = c.ShouldBindJSON(&password); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// retrieve user
	var user models.User
	if user, err = handlers.Handler.UserGetOneBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	// update user password
	if err := handlers.Handler.UserUpdatePassword(c.GetString(constants.UserSubjectKey), password.Old, password.New); err != nil {
		if err == errors.ErrAuthIncorrectCredentials {
			handlers.Handler.LogInsertUser(user, false, datatransfers.LogDetail{
				Scope:  constants.LogScopeUserPassword,
				Detail: utils.ParseUserAgent(c),
			})
			c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Code: err.Error(), Error: "incorrect old_password"})
		} else {
			c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "failed updating user password"})
		}
		return
	}
	handlers.Handler.LogInsertUser(user, true, datatransfers.LogDetail{Scope: constants.LogScopeUserPassword})
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Update user superuser status
// @Tags user
// @Param user body datatransfers.UserUpdateSuperuser true "User update superuser info"
// @Success 200 "OK"
// @Router /api/v1/user/superuser [PUT]
func PUTUserSuperuser(c *gin.Context) {
	var err error
	// fetch superuser info
	var superuser datatransfers.UserUpdateSuperuser
	if err = c.ShouldBindJSON(&superuser); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// prevent changing own status
	if c.GetString(constants.UserSubjectKey) == superuser.Subject {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "cannot change own superuser status"})
		return
	}
	// retrieve users
	var user, target models.User
	if user, err = handlers.Handler.UserGetOneBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	if target, err = handlers.Handler.UserGetOneBySubject(superuser.Subject); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	// clear session
	var activeSessions []*datatransfers.SessionInfo
	if activeSessions, err = handlers.Handler.SessionGetAllActive(superuser.Subject); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "cannot retrieve active sessions"})
		return
	}
	for _, session := range activeSessions {
		handlers.Handler.SessionRevoke(session.SessionID)
	}
	// update status
	if err = handlers.Handler.UserUpdateSuperuser(superuser.Subject, superuser.Superuser); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating user"})
		return
	}
	handlers.Handler.LogInsertUser(user, superuser.Superuser, datatransfers.LogDetail{
		Scope:  constants.LogScopeUserSuperuser,
		Detail: datatransfers.UserInfo{Username: target.Username},
	})
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
	if err := handlers.Handler.AuthVerify(verify.Token); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "failed verifying user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Resend verification email
// @Tags user
// @Param user body datatransfers.UserResend true "User resend info"
// @Success 200 "OK"
// @Router /api/v1/user/resend [POST]
func POSTResend(c *gin.Context) {
	var err error
	// fetch verification info
	var resend datatransfers.UserResend
	if err = c.ShouldBindJSON(&resend); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// get user
	var user models.User
	if user, err = handlers.Handler.UserGetOneByEmail(resend.Email); err != nil {
		c.JSON(http.StatusOK, datatransfers.APIResponse{}) // silent request drop
		return
	}
	// resend email
	if !user.EmailVerified { // silent request drop
		if err := handlers.Handler.MailerSendEmailVerification(user); err != nil {
			c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "failed verifying user"})
			return
		}
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
