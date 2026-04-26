package main

import (
	"log"
	"net/http"
	"time"

	core_constants "buckly-ms/core/constants"
	"buckly-ms/gateway/di"
	_ "buckly-ms/gateway/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	container := di.BuildContainer()

	err := container.Invoke(func(p di.Params) {
		var logger *zap.Logger
		if p.Config.Environment == core_constants.DEVELOPMENT {
			var e error
			logger, e = zap.NewDevelopment()
			if e != nil {
				log.Fatalf("Failed to initialize logger: %v", e)
			}
			logger.Info("Running in development mode")
			gin.SetMode(gin.DebugMode)
		} else {
			var e error
			logger, e = zap.NewProduction()
			if e != nil {
				log.Fatalf("Failed to initialize logger: %v", e)
			}
			logger.Info("Running in production mode")
			gin.SetMode(gin.ReleaseMode)
		}

		r := gin.Default()

		for _, handler := range p.Handler {
			handler.RegisterRoutes(r)
		}

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		srv := &http.Server{
			Addr:         ":" + p.Config.HTTPPort,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		logger.Info("Starting gateway service", zap.Any("config: ", p.Config))

		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Failed to run gateway: %v", err)
		}

	})

	if err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}

}

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns status ok if the service is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
