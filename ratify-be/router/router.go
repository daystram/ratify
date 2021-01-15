package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/daystram/ratify/ratify-be/controllers/middleware"
	"github.com/daystram/ratify/ratify-be/controllers/oauth"
	v1 "github.com/daystram/ratify/ratify-be/controllers/v1"
	_ "github.com/daystram/ratify/ratify-be/docs" // ininitliase SwaggerUI
	"github.com/daystram/ratify/ratify-be/utils"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))
	apiV1 := router.Group("/api/v1") // internal/dashboard APIs
	apiV1.Use(
		middleware.AuthMiddleware,
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
			user.PUT("/password", utils.AuthOnly, v1.PUTUserPassword)
			user.POST("/verify", v1.POSTVerify)
			user.POST("/resend", v1.POSTResend)
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
		mfa := apiV1.Group("/mfa")
		{
			mfa.POST("/enable", utils.AuthOnly, v1.POSTEnableTOTP)
			mfa.POST("/confirm", utils.AuthOnly, v1.POSTConfirmTOTP)
			mfa.POST("/disable", utils.AuthOnly, v1.POSTDisableTOTP)
		}
		log := apiV1.Group("/log")
		{
			log.GET("/user_activity", utils.AuthOnly, v1.GETActivityLog)
			log.GET("/admin_activity", utils.AuthOnly, utils.SuperuserOnly, v1.GETAdminLog)
		}
	}
	oauthV1 := router.Group("/oauth") // OAuth
	oauthV1.Use(
		middleware.CORSMiddleware,
	)
	{
		oauthV1.POST("/authorize", oauth.POSTAuthorize)
		oauthV1.POST("/token", oauth.POSTToken)
		oauthV1.POST("/introspect", oauth.POSTIntrospect)
		oauthV1.GET("/userinfo", middleware.AuthMiddleware, utils.AuthOnly, oauth.GETUserInfo)
		oauthV1.POST("/logout", oauth.POSTLogout)
	}
	return
}
