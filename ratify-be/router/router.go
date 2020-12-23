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
	v1route := router.Group("/api/v1") // internal/dashboard APIs
	v1route.Use(
		middleware.CORSMiddleware,
		middleware.AuthMiddleware,
	)
	{
		auth := v1route.Group("/auth")
		{
			auth.POST("/login", v1.POSTLogin)
			auth.POST("/signup", v1.POSTRegister)
		}
		user := v1route.Group("/user")
		{
			user.GET("/:username", utils.AuthOnly, v1.GETUser)
			user.PUT("/", utils.AuthOnly, v1.PUTUser)
		}
		application := v1route.Group("/application")
		{
			application.GET("/", utils.AuthOnly, utils.SuperuserOnly, v1.GETOwnedApplications)
			application.GET("/:client_id", utils.AuthOnly, utils.SuperuserOnly, v1.GETOneApplication)
			application.POST("/", utils.AuthOnly, utils.SuperuserOnly, v1.POSTApplication)
			application.PUT("/:client_id", utils.AuthOnly, utils.SuperuserOnly, v1.PUTApplication)
		}
	}
	return
}
