package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/handlers"
)

// @Summary Check unique form field
// @Tags form
// @Param uniqueRequest body datatransfers.UniqueCheckRequest true "Unique query"
// @Success 200 "OK"
// @Router /api/v1/form/unique [POST]
func POSTUniqueCheck(c *gin.Context) {
	var err error
	var uniqueRequest datatransfers.UniqueCheckRequest
	if err = c.ShouldBindJSON(&uniqueRequest); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: err.Error()})
		return
	}
	unique := false
	switch uniqueRequest.Field {
	case "user:username":
		_, err = handlers.Handler.RetrieveUserByUsername(uniqueRequest.Value)
		unique = err != nil
	case "user:email":
		_, err = handlers.Handler.RetrieveUserByEmail(uniqueRequest.Value)
		unique = err != nil
	default:
		c.JSON(http.StatusBadRequest, datatransfers.APIResponse{Error: fmt.Sprintf("unsupported field %s", uniqueRequest.Field)})
		return
	}
	c.JSON(http.StatusOK, datatransfers.APIResponse{Data: unique})
	return
}