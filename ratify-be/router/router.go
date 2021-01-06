package router

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/daystram/ratify/ratify-be/controllers/middleware"
	"github.com/daystram/ratify/ratify-be/controllers/oauth"
	"github.com/daystram/ratify/ratify-be/controllers/v1"
	_ "github.com/daystram/ratify/ratify-be/docs"
	"github.com/daystram/ratify/ratify-be/utils"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))
	apiV1 := router.Group("/api/v1") // internal/dashboard APIs
	apiV1.Use(
		middleware.CORSMiddleware,
		middleware.AuthMiddleware,
		utils.LogIP,
	)
	{
		form := apiV1.Group("/form")
		{
			form.POST("/unique", v1.POSTUniqueCheck)
		}
		user := apiV1.Group("/user")
		{
			user.GET("/", utils.AuthOnly, v1.GETUser)
			user.POST("/", v1.POSTRegister)
			user.PUT("/", utils.AuthOnly, v1.PUTUser)
		}
		application := apiV1.Group("/application")
		{
			application.GET("/", utils.AuthOnly, utils.SuperuserOnly, v1.GETApplicationList)
			application.GET("/:client_id", v1.GETOneApplicationDetail)
			application.POST("/", utils.AuthOnly, utils.SuperuserOnly, v1.POSTApplication)
			application.PUT("/:client_id", utils.AuthOnly, utils.SuperuserOnly, v1.PUTApplication)
			application.PUT("/:client_id/revoke", utils.AuthOnly, utils.SuperuserOnly, v1.PUTApplicationRevokeSecret)
			application.DELETE("/:client_id", utils.AuthOnly, utils.SuperuserOnly, v1.DELETEApplication)
		}
	}
	oauthV1 := router.Group("/oauth") // oauth
	oauthV1.Use(
		middleware.CORSMiddleware,
	)
	{
		oauthV1.POST("/authorize", oauth.POSTAuthorize)
		oauthV1.POST("/token", oauth.POSTToken)
		oauthV1.POST("/introspect", oauth.POSTIntrospect)
		oauthV1.POST("/logout", middleware.AuthMiddleware, utils.AuthOnly, oauth.POSTLogout)
	}
	return
}
