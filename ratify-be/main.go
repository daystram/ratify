package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daystram/ratify/ratify-be/config"
	"github.com/daystram/ratify/ratify-be/docs"
	"github.com/daystram/ratify/ratify-be/handlers"
	"github.com/daystram/ratify/ratify-be/router"
)

// @title Ratify
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @scheme Bearer

func init() {
	config.InitializeAppConfig()
	if !config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	docs.SwaggerInfo.Host = config.AppConfig.Hostname
	handlers.InitializeHandler()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        router.InitializeRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
