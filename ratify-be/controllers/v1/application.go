package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	// fetch application info
	var application models.Application
	application.ClientID = strings.TrimPrefix(c.Param("client_id"), "/") // trim due to router catch-all
	if application, err = handlers.Handler.RetrieveApplication(application.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// check superuser
	if !c.GetBool(constants.IsSuperuserKey) {
		var next int
		if sessionCookie, err := c.Request.Cookie("test"); err == nil {
			// user still signed in, check cookie
			next, _ = strconv.Atoi(sessionCookie.Value)
		} else {
			log.Println(err)
		}
		c.SetCookie("test", fmt.Sprintf("%d", next + 1), 120, "/api/v1/application", "localhost", false, true)
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
func GETApplicationList(c *gin.Context) {
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
	var clientSecret string
	if applicationInfo.ClientID, clientSecret, err = handlers.Handler.RegisterApplication(applicationInfo, c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: gin.H{"client_id": applicationInfo.ClientID, "client_secret": clientSecret}})
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
	var application models.Application
	application.ClientID = strings.TrimPrefix(c.Param("client_id"), "/")
	if application, err = handlers.Handler.RetrieveApplication(application.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// checked locked flag
	if application.Locked &&
		(applicationInfo.LoginURL != application.LoginURL ||
			applicationInfo.CallbackURL != application.CallbackURL ||
			applicationInfo.LogoutURL != application.LogoutURL) {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "application is locked"})
		return
	}
	// update application
	applicationInfo.ClientID = application.ClientID
	if err = handlers.Handler.UpdateApplication(applicationInfo); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed updating application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Delete application
// @Tags application
// @Security BearerAuth
// @Param client_id path string true "Client ID"
// @Success 200 "OK"
// @Router /api/v1/application/{client_id} [DELETE]
func DELETEApplication(c *gin.Context) {
	var err error
	// fetch application info
	var application models.Application
	application.ClientID = strings.TrimPrefix(c.Param("client_id"), "/")
	if application, err = handlers.Handler.RetrieveApplication(application.ClientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// checked locked flag
	if application.Locked {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: "application is locked"})
		return
	}
	// delete application
	if err = handlers.Handler.DeleteApplication(application.ClientID); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed deleting application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{})
	return
}

// @Summary Revoke application secret
// @Tags application
// @Security BearerAuth
// @Param client_id path string true "Client ID"
// @Success 200 "OK"
// @Router /api/v1/application/{client_id}/revoke [PUT]
func PUTApplicationRevokeSecret(c *gin.Context) {
	var err error
	// fetch application info
	clientID := strings.TrimPrefix(c.Param("client_id"), "/")
	if _, err = handlers.Handler.RetrieveApplication(clientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.APIResponse{Error: "application not found"})
		return
	}
	// renew application client_secret
	var clientSecret string
	if clientSecret, err = handlers.Handler.RenewApplicationClientSecret(clientID); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.APIResponse{Error: "failed renewing application client_secret"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: gin.H{"client_secret": clientSecret}})
	return
}
