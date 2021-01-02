package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
)

func AuthOnly(c *gin.Context) {
	if !c.GetBool(constants.IsAuthenticatedKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: "invalid_token", Error: "user not authenticated"})
	}
}

func SuperuserOnly(c *gin.Context) {
	if !c.GetBool(constants.IsSuperuserKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: "unauthorized", Error: "access unauthorized"})
	}
}
