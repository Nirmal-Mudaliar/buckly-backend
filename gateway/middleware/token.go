package middleware

import (
	core_constants "buckly-ms/core/constants"
	core_utils "buckly-ms/core/utils"
	"buckly-ms/gateway/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := core_utils.GetLoggerFromContext(c)
		log.Info("Token middleware executed")
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Success: false,
				Message: "Unauthorized: No token provided",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// For renew-access-token endpoint, allow expired tokens but validate signature
		isRefreshEndpoint := strings.Contains(c.Request.URL.Path, "auth/renew-access-token")

		var claims *core_utils.Claims
		var err error

		if isRefreshEndpoint {
			// Validate token signature and claims, but allow expired tokens
			claims, err = core_utils.ValidateTokenIgnoreExpiration(tokenString, jwtSecret)
		} else {
			// Strict validation - reject expired tokens
			claims, err = core_utils.ValidateToken(tokenString, jwtSecret)
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(core_constants.JWT_TOKEN_CLAIM, claims)
		c.Next()
	}

}
