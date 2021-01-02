package router

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/daystram/ratify/ratify-be/controllers/middleware"
	"github.com/daystram/ratify/ratify-be/controllers/v1"
	_ "github.com/daystram/ratify/ratify-be/docs"
	"github.com/daystram/ratify/ratify-be/utils"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))
	apiv1 := router.Group("/api/v1") // internal/dashboard APIs
	apiv1.Use(
		middleware.CORSMiddleware,
		middleware.AuthMiddleware,
	)
	{
		form := apiv1.Group("/form")
		{
			form.POST("/unique", v1.POSTUniqueCheck)
		}
		user := apiv1.Group("/user")
		{
			user.GET("/", utils.AuthOnly, v1.GETUser)
			user.POST("/", v1.POSTRegister)
			user.PUT("/", utils.AuthOnly, v1.PUTUser)
		}
		application := apiv1.Group("/application")
		{
			application.GET("/", utils.AuthOnly, utils.SuperuserOnly, v1.GETApplicationList)
			application.GET("/:client_id", v1.GETOneApplicationDetail)
			application.POST("/", utils.AuthOnly, utils.SuperuserOnly, v1.POSTApplication)
			application.PUT("/:client_id", utils.AuthOnly, utils.SuperuserOnly, v1.PUTApplication)
			application.PUT("/:client_id/revoke", utils.AuthOnly, utils.SuperuserOnly, v1.PUTApplicationRevokeSecret)
			application.DELETE("/:client_id", utils.AuthOnly, utils.SuperuserOnly, v1.DELETEApplication)
		}
	}
	oauth := router.Group("/oauth") // oauth
	oauth.Use(
		middleware.CORSMiddleware,
	)
	{
		oauth.POST("/authorize", v1.POSTAuthorize)
		oauth.POST("/token", v1.POSTToken)
		oauth.POST("/introspect", v1.POSTIntrospect)
	}
	return
}
