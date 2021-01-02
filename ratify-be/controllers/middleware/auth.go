package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/models"
)

func AuthMiddleware(c *gin.Context) {
	accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if accessToken == "" {
		c.Set(constants.IsAuthenticatedKey, false)
		c.Next()
		return
	}
	var err error
	var tokenInfo datatransfers.TokenIntrospection
	if tokenInfo, err = handlers.Handler.IntrospectAccessToken(accessToken); err != nil || !tokenInfo.Active {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.APIResponse{Code: "invalid_token", Error: "invalid access_token"})
		return
	}
	var user models.User
	if user, err = handlers.Handler.RetrieveUserBySubject(tokenInfo.Subject); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.APIResponse{Code:"invalid_token", Error: err.Error()})
		return
	}
	c.Set(constants.IsAuthenticatedKey, true)
	c.Set(constants.UserSubjectKey, user.Subject)
	c.Set(constants.IsSuperuserKey, user.Superuser)
	c.Next()
}
