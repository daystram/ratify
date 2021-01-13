package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"

	"github.com/daystram/ratify/ratify-be/datatransfers"
)

func ParseUserAgent(c *gin.Context) datatransfers.UserAgent {
	ua := user_agent.New(c.Request.UserAgent())
	browser, _ := ua.Browser()
	return datatransfers.UserAgent{
		IP:      c.ClientIP(),
		Browser: browser,
		OS:      ua.OS(),
	}
}
