package main

import (
	"log"
	"net/http"
	"time"

	_ "buckly-ms/gateway/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	r := gin.Default()
	r.GET("/health", health)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:         ":" + "8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to run gateway: %v", err)
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
