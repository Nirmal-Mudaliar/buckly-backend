package auth

import (
	"buckly-ms/gateway/api/handlers/auth/dto"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	"buckly-ms/gateway/models"
	auth_gen "buckly-ms/proto/auth-gen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Config            *config.GatewayConfig
	AuthServiceClient auth_gen.AuthServiceClient
}

func NewAuthHandler(
	config *config.GatewayConfig,
	authServiceClient auth_gen.AuthServiceClient,
) contracts.RouteRegistrar {
	return &AuthHandler{
		Config:            config,
		AuthServiceClient: authServiceClient,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	secured := r.Group("/api/v1/auth")
	secured.POST("/signup", h.SignUp)
	secured.POST("/verify-otp", h.VerifyOTP)
	secured.POST("/renew-access-token", h.RenewAccessToken)
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
		Email:       payload.Email,
		Password:    payload.Password,
		PhoneNo:     payload.PhoneNo,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		DateOfBirth: payload.DateOfBirth,
		Gender:      payload.Gender,
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

// Verify OTP godoc
// @Summary      Verify OTP
// @Description  Verifies the OTP sent to the user's phone number during sign-up and returns access and refresh tokens upon successful verification
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.VerifyOTPRequest  true  "Verify OTP details"
// @Success      200  {object}  models.ApiResponse{data=dto.VerifyOTPResponse}
// @Router       /api/v1/auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var payload dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Payload",
		})
		return
	}

	if payload.PhoneNo == "" || payload.OTP == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Phone number and OTP are required",
		})
		return
	}

	resp, err := h.AuthServiceClient.VerifyOTP(ctx, &auth_gen.VerifyOTPRequest{
		PhoneNo: payload.PhoneNo,
		Otp:     payload.OTP,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occured in OTP verification",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "OTP verification successful",
		Data: dto.VerifyOTPResponse{
			User: &dto.User{
				Id:              resp.User.Id,
				FirstName:       resp.User.FirstName,
				LastName:        resp.User.LastName,
				Email:           resp.User.Email,
				PhoneNo:         resp.User.PhoneNo,
				DateOfBirth:     resp.User.DateOfBirth,
				Gender:          resp.User.Gender,
				Bio:             resp.User.Bio,
				ProfilePhotoUrl: resp.User.ProfilePhotoUrl,
				HomeCountryId:   resp.User.HomeCountryId,
				HomeStateId:     resp.User.HomeStateId,
				HomeCityId:      resp.User.HomeCityId,
				IsPhoneVerified: resp.User.IsPhoneVerified,
				TrustScore:      resp.User.TrustScore,
				Status:          resp.User.Status,
				InsertTs:        resp.User.InsertTs.AsTime().String(),
				ModifiedTs:      resp.User.ModifiedTs.AsTime().String(),
			},
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
		},
	})
}

// Renew Access Token godoc
// @Summary      Renew Access Token
// @Description  Renews the access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.RenewAccessTokenRequest  true  "Renew access token details"
// @Success      200  {object}  models.ApiResponse{data=dto.RenewAccessTokenResponse}
// @Router       /api/v1/auth/renew-access-token [post]
func (h *AuthHandler) RenewAccessToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var payload dto.RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Payload",
		})
		return
	}

	resp, err := h.AuthServiceClient.RenewAccessToken(
		ctx,
		&auth_gen.RenewAccessTokenRequest{
			RefreshToken: payload.RefreshToken,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occured in renewing access token",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Access token renewed successfully",
		Data: dto.RenewAccessTokenResponse{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
		},
	})
}
