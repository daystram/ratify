package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

// @Summary Get all activity log for current user
// @Tags log
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/log/user_activity [GET]
func GETActivityLog(c *gin.Context) {
	var err error
	// get logs
	var logs []models.Log
	if logs, err = handlers.Handler.RetrieveActivityLogs(c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "cannot retrieve activity logs"})
		return
	}
	logsResponse := make([]datatransfers.LogInfo, 0)
	for _, entry := range logs {
		logsResponse = append(logsResponse, datatransfers.LogInfo{
			ClientID:    entry.Application.ClientID,
			Severity:    entry.Severity,
			Description: entry.Description,
			CreatedAt:   entry.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: logsResponse})
	return
}

// @Summary Get all admin application log
// @Tags log
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/log/admin_activity [GET]
func GETAdminLog(c *gin.Context) {
	var err error
	// get logs
	var logs []models.Log
	if logs, err = handlers.Handler.RetrieveAdminLogs(); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "cannot retrieve admin logs"})
		return
	}
	log.Println(logs)
	logsResponse := make([]datatransfers.LogInfo, 0)
	for _, entry := range logs {
		logsResponse = append(logsResponse, datatransfers.LogInfo{
			Username:        entry.User.Username,
			ApplicationName: entry.Application.Name,
			ClientID:        entry.Application.ClientID,
			Severity:        entry.Severity,
			Description:     entry.Description,
			CreatedAt:       entry.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: logsResponse})
	return
}
