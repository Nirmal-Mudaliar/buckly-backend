package auth

import (
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	"buckly-ms/gateway/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Config *config.GatewayConfig
}

func NewAuthHandler(config *config.GatewayConfig) contracts.RouteRegistrar {
	return &AuthHandler{
		Config: config,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	secured := r.Group("/api/v1/auth")
	secured.GET("/signup", h.SignUp)
}

// Sign Up godoc
// @Summary      Signs up the user
// @Description  Returns status ok if the service is running
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.ApiResponse
// @Router       /api/v1/auth/signup [get]
func (h *AuthHandler) SignUp(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Signup successful",
	})
}
