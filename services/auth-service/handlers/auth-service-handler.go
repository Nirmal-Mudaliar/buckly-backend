package handlers

import (
	"buckly-ms/core/utils"
	auth_gen "buckly-ms/proto/auth-gen"
	database_gen "buckly-ms/proto/database-gen"
	auth_utils "buckly-ms/services/auth-service/utils"
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (ass *AuthServiceServer) SignUp(ctx context.Context, req *auth_gen.SignUpRequest) (*auth_gen.SignUpResponse, error) {
	logger := utils.GetLoggerFromContext(ctx)

	existingUserResp, err := ass.DatabaseServiceClient.GetUserByEmail(ctx, &database_gen.GetUserByEmailRequest{
		Email: req.Email,
	})

	if err == nil && existingUserResp.User != nil {
		if existingUserResp.User.IsPhoneVerified {
			return nil, status.Errorf(codes.AlreadyExists, "Email already registered")
		} else {
			logger.Info("User already exist with unverified phone no")
			_, err := ass.DatabaseServiceClient.DeleteUnverifiedUserByUserId(ctx, &database_gen.DeleteUnverifiedUserByUserIdRequest{
				Id: existingUserResp.User.Id,
			})
			if err != nil {
				logger.Error("Failed to delete existing user who's phone no is not verified", zap.Error(err))
				return nil, status.Errorf(codes.Internal, "Failed to delete existing user")
			}
			logger.Info("User deleted successfully")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to process password")
	}

	_, err = ass.DatabaseServiceClient.CreateUser(ctx, &database_gen.CreateUserRequest{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNo:      req.PhoneNo,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		logger.Error("Failed to create user: ", zap.Error(err))
		return nil, err
	}

	if err := ass.TwilioClient.SendOTP(req.PhoneNo); err != nil {
		logger.Error("Failed to send OTP", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to send OTP: %v", err)
	}

	return &auth_gen.SignUpResponse{
		Message: "OTP sent to you phone",
	}, nil
}

func (ass *AuthServiceServer) VerifyOTP(ctx context.Context, req *auth_gen.VerifyOTPRequest) (*auth_gen.VerifyOTPResponse, error) {
	logger := utils.GetLoggerFromContext(ctx)

	isApproved, err := ass.TwilioClient.VerifyOTP(req.PhoneNo, req.Otp)
	if err != nil {
		logger.Error("Failed to verify OTP", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to verify OTP: %v", err)
	}

	if !isApproved {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid OTP")
	}

	userResp, err := ass.DatabaseServiceClient.GetUserByPhoneNo(ctx, &database_gen.GetUserByPhoneNoRequest{
		PhoneNo: req.PhoneNo,
	})

	if err != nil {
		logger.Error("Failed to get user by phone no", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get user by phone no")
	}

	updatePhoneNoVerifiedResp, err := ass.DatabaseServiceClient.UpdatePhoneNoVerified(ctx, &database_gen.UpdatePhoneNoVerifiedRequest{
		PhoneNo: req.PhoneNo,
	})

	if err != nil {
		logger.Error("Failed to update phone no verified", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to update phone no verified")
	}

	accessToken, err := utils.GenerateAccessToken(userResp.User.Id, userResp.User.Email, req.PhoneNo, true, ass.JWTSecret)
	if err != nil {
		logger.Error("Failed to generate access token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(userResp.User.Id, ass.JWTSecret)
	if err != nil {
		logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to generate refresh token")
	}

	_, err = ass.DatabaseServiceClient.StoreRefreshToken(ctx, &database_gen.StoreRefreshTokenRequest{
		UserId:    userResp.User.Id,
		Token:     refreshToken,
		ExpiresAt: timestamppb.New(time.Now().Add(30 * 24 * time.Hour)),
	})

	if err != nil {
		logger.Error("Failed to store refresh token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to store refresh token")
	}

	return &auth_gen.VerifyOTPResponse{
		User:         auth_utils.DatabaseUserToAuthUser(updatePhoneNoVerifiedResp.User),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
