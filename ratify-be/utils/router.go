package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/constants"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/datatransfers"
)

func AuthOnly(c *gin.Context) {
	if !c.GetBool(constants.IsAuthenticatedKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: "user not authenticated"})
	}
}
