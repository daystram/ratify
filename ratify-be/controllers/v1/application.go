package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

// @Summary Get application details
// @Description If unauthorized, only application name is returned
// @Tags application
// @Security BearerAuth
// @Param client_id path string true "Client ID"
// @Success 200 "OK"
// @Router /api/v1/application/{client_id} [GET]
func GETOneApplicationDetail(c *gin.Context) {
	var err error
	// fetch clientID
	clientID := strings.TrimPrefix(c.Param("client_id"), "/") // trim due to router catch-all
	// get application
	var application models.Application
	if application, err = handlers.Handler.RetrieveApplication(clientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// check superuser
	if !c.GetBool(constants.IsSuperuserKey) {
		c.JSON(http.StatusOK, datatransfers.APIResponse{Data: datatransfers.ApplicationInfo{
			Name: application.Name,
		}})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: datatransfers.ApplicationInfo{
		ClientID:    application.ClientID,
		Name:        application.Name,
		Description: application.Description,
		LoginURL:    application.LoginURL,
		CallbackURL: application.CallbackURL,
		LogoutURL:   application.LogoutURL,
		Metadata:    application.Metadata,
		Locked:      application.Locked,
		CreatedAt:   application.CreatedAt,
	}})
	return
}

// @Summary Get all applications
// @Tags application
// @Security BearerAuth
// @Success 200 "OK"
// @Router /api/v1/application [GET]
func GETApplications(c *gin.Context) {
	var err error
	// get all owned applications
	var applications []models.Application
	if applications, err = handlers.Handler.RetrieveAllApplications(); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "cannot retrieve applications"})
		return
	}
	var applicationsResponse []datatransfers.ApplicationInfo
	for _, application := range applications {
		applicationsResponse = append(applicationsResponse, datatransfers.ApplicationInfo{
			ClientID:    application.ClientID,
			Name:        application.Name,
			Description: application.Description,
			CreatedAt:   application.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: applicationsResponse})
	return
}

// @Summary Create application
// @Tags application
// @Security BearerAuth
// @Param application body datatransfers.ApplicationInfo true "Application info"
// @Success 200 "OK"
// @Router /api/v1/application [POST]
func POSTApplication(c *gin.Context) {
	var err error
	// fetch application info
	var applicationInfo datatransfers.ApplicationInfo
	if err = c.ShouldBindJSON(&applicationInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	// register application
	if applicationInfo.ClientID, err = handlers.Handler.RegisterApplication(applicationInfo, c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: gin.H{"client_id": applicationInfo.ClientID}})
	return
}

// @Summary Update application
// @Tags application
// @Security BearerAuth
// @Param client_id path string true "Client ID"
// @Param application body datatransfers.ApplicationInfo true "Application info"
// @Success 200 "OK"
// @Router /api/v1/application/{client_id} [PUT]
func PUTApplication(c *gin.Context) {
	var err error
	// fetch application info
	var applicationInfo datatransfers.ApplicationInfo
	if err = c.ShouldBindJSON(&applicationInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	applicationInfo.ClientID = c.Param("client_id")
	if _, err = handlers.Handler.RetrieveApplication(applicationInfo.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// update application
	if err = handlers.Handler.UpdateApplication(applicationInfo); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}
