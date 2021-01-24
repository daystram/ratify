package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

// @Summary Get dashboard info
// @Tags dashboard
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/dashboard [GET]
func GETDashboardInfo(c *gin.Context) {
	var err error
	// retrieve user
	var user models.User
	if user, err = handlers.Handler.UserGetOneBySubject(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "user not found"})
		return
	}
	// retrieve all sessions
	var activeSessions []*datatransfers.SessionInfo
	if activeSessions, err = handlers.Handler.SessionGetAllActive(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "cannot retrieve active sessions"})
		return
	}
	dashboardInfo := datatransfers.DashboardInfo{
		SignInCount:  user.LoginCount,
		LastSignIn:   user.LastLogin,
		SessionCount: len(activeSessions),
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: dashboardInfo})
	return
}
