package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/constants"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/datatransfers"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/handlers"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/models"
)

func GETApplication(c *gin.Context) {
	var err error
	// fetch clientID
	clientID := strings.TrimPrefix(c.Param("client_id"), "/") // trim due to router catch-all
	if len(clientID) == 0 {
		// get all owned applications
		var applications []models.Application
		if applications, err = handlers.Handler.RetrieveOwnedApplications(c.GetString(constants.UserSubjectKey)); err != nil {
			c.JSON(http.StatusNotFound, datatransfers.Response{Error: "cannot retrieve applications"})
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
		c.JSON(http.StatusOK, datatransfers.Response{Data: applicationsResponse})
	} else {
		// get one application
		var application models.Application
		if application, err = handlers.Handler.RetrieveApplication(clientID); err != nil {
			c.JSON(http.StatusNotFound, datatransfers.Response{Error: "application not found"})
			return
		}
		// check ownership
		if application.Owner.Subject != c.GetString(constants.UserSubjectKey) {
			c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "access to resource unauthorized"})
			return
		}
		c.JSON(http.StatusOK, datatransfers.Response{Data: datatransfers.ApplicationInfo{
			ClientID:     application.ClientID,
			ClientSecret: application.ClientSecret,
			Name:         application.Name,
			Description:  application.Description,
			LoginURL:     application.LoginURL,
			CallbackURL:  application.CallbackURL,
			LogoutURL:    application.LogoutURL,
			Metadata:     application.Metadata,
			CreatedAt:    application.CreatedAt,
		}})
	}
	return
}

func POSTApplication(c *gin.Context) {
	var err error
	// fetch application info
	var applicationInfo datatransfers.ApplicationInfo
	if err = c.ShouldBindJSON(&applicationInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}
	// register application
	if applicationInfo.ClientID, err = handlers.Handler.RegisterApplication(applicationInfo, c.GetString(constants.UserSubjectKey)); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.Response{Error: "failed updating application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.Response{Data: gin.H{"client_id": applicationInfo.ClientID}})
	return
}

func PUTApplication(c *gin.Context) {
	var err error
	// fetch application info
	clientID := c.Param("client_id")
	var applicationInfo datatransfers.ApplicationInfo
	if err = c.ShouldBindJSON(&applicationInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}
	// check ownership
	var application models.Application
	if application, err = handlers.Handler.RetrieveApplication(clientID); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "application not found"})
		return
	}
	if application.Owner.Subject != c.GetString(constants.UserSubjectKey) {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "access to resource unauthorized"})
		return
	}
	// update application
	if err = handlers.Handler.UpdateApplication(applicationInfo); err != nil {
		c.JSON(http.StatusInternalServerError, datatransfers.Response{Error: "failed updating application"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.Response{})
	return
}
