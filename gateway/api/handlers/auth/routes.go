package auth

import (
	"buckly-ms/gateway/api/handlers/auth/dto"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	"buckly-ms/gateway/models"
	auth_gen "buckly-ms/proto/auth-gen"
	database_gen "buckly-ms/proto/database-gen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Config                *config.GatewayConfig
	DatabaseServiceClient database_gen.DatabaseServiceClient
	AuthServiceClient     auth_gen.AuthServiceClient
}

func NewAuthHandler(
	config *config.GatewayConfig,
	databaseServiceClient database_gen.DatabaseServiceClient,
	authServiceClient auth_gen.AuthServiceClient,
) contracts.RouteRegistrar {
	return &AuthHandler{
		Config:                config,
		DatabaseServiceClient: databaseServiceClient,
		AuthServiceClient:     authServiceClient,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	secured := r.Group("/api/v1/auth")
	secured.POST("/signup", h.SignUp)
}

// Sign Up godoc
// @Summary      Register a new user
// @Description  Creates a new user account with email, password, and profile information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.SignUpRequest  true  "Sign up details"
// @Success      200  {object}  models.ApiResponse{data=dto.SignUpResponse}
// @Router       /api/v1/auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var payload dto.SignUpRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Payload",
		})
	}

	resp, err := h.AuthServiceClient.SignUp(ctx, &auth_gen.SignUpRequest{
		Email:      payload.Email,
		Password:   payload.Password,
		PhoneNo:    payload.PhoneNo,
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		DateOfBith: payload.DateOfBirth,
		Gender:     payload.Gender,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occured in Sign Up",
		})
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Signup successful",
		Data: dto.SignUpResponse{
			Message: resp.Message,
		},
	})
}
