package handlers

import (
	"buckly-ms/core/utils"
	auth_gen "buckly-ms/proto/auth-gen"
	database_gen "buckly-ms/proto/database-gen"
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ass *AuthServiceServer) SignUp(ctx context.Context, req *auth_gen.SignUpRequest) (*auth_gen.SignUpResponse, error) {
	logger := utils.GetLoggerFromContext(ctx)

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
		DateOfBirth:  req.DateOfBith,
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
