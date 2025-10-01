package app

import (
	"github.com/gin-contrib/cors"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"oolio.com/kart/configs"
	"oolio.com/kart/constants"
	"oolio.com/kart/controllers"
	"time"
)

func initializeRoutes() *gin.Engine {
	if configs.ReleaseEnv == constants.ProdMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(cors.Default())

	router.Use(ginZap.RecoveryWithZap(configs.Logger, true))
	router.Use(ginZap.Ginzap(configs.Logger, time.RFC3339, false))

	kartRouter := router.Group("/kart/v1")

	kartRouter.GET("/health", controllers.HealthCheckController.HealthCheck)

	return router
}
