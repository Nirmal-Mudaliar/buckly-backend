package auth

import (
	"buckly-ms/gateway/api/handlers/auth/dto"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	"buckly-ms/gateway/models"
	database_gen "buckly-ms/proto/database-gen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	Config                *config.GatewayConfig
	DatabaseServiceClient database_gen.DatabaseServiceClient
}

func NewAuthHandler(config *config.GatewayConfig, databaseServiceClient database_gen.DatabaseServiceClient) contracts.RouteRegistrar {
	return &AuthHandler{
		Config:                config,
		DatabaseServiceClient: databaseServiceClient,
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	usersResp, err := h.DatabaseServiceClient.GetAllUsers(ctx, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Failed to fetch users from database service",
		})
		return
	}

	users := make([]dto.Users, 0, len(usersResp.Users))
	for _, u := range usersResp.Users {
		users = append(users, dto.Users{
			Id:              u.Id,
			FirstName:       u.FirstName,
			LastName:        u.LastName,
			Email:           u.Email,
			PhoneNo:         u.PhoneNo,
			DateOfBirth:     u.DateOfBirth,
			Gender:          u.Gender,
			Bio:             u.Bio,
			ProfilePhotoUrl: u.ProfilePhotoUrl,
			HomeCountryId:   u.HomeCountryId,
			HomeStateId:     u.HomeStateId,
			HomeCityId:      u.HomeCityId,
			IsPhoneVerified: u.IsPhoneVerified,
			TrustScore:      u.TrustScore,
			Status:          u.Status,
			InsertTs:        u.InsertTs.AsTime().Format(time.RFC3339),
			ModifiedTs:      u.ModifiedTs.AsTime().Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Signup successful",
		Data: dto.SignUpResponse{
			Users: users,
		},
	})
}
