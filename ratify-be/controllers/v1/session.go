package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
)

// @Summary Get all active sessions of current user
// @Tags session
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/session [GET]
func GETSessionActive(c *gin.Context) {
	var err error
	// retrieve all sessions
	var activeSessions []*datatransfers.SessionInfo
	if activeSessions, err = handlers.Handler.SessionGetAllActive(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "cannot retrieve active sessions"})
		return
	}
	// check current session
	var sessionID = c.GetString(constants.SessionIDKey)
	for _, session := range activeSessions {
		session.Current = (session.SessionID == sessionID)
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: activeSessions})
	return
}

// @Summary Revoke session
// @Tags session
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/session [POST]
func POSTSessionRevoke(c *gin.Context) {
	var err error
	var session datatransfers.SessionInfo
	if err = c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	if err = handlers.Handler.SessionRevoke(session.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed clearing session"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
