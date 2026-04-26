package handlers

import (
	database_gen "buckly-ms/proto/database-gen"
	db "buckly-ms/services/database-service/db/generated"
	"buckly-ms/services/database-service/utils"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (dss *DatabaseServiceServer) GetAllUsers(ctx context.Context, empty *emptypb.Empty) (*database_gen.GetAllUsersResponse, error) {
	rows, err := dss.Queries.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error fetching users from database: %v", err)
	}
	users := make([]*database_gen.User, 0, len(rows))
	for _, r := range rows {
		users = append(users, &database_gen.User{
			Id:              r.ID,
			FirstName:       r.FirstName,
			LastName:        r.LastName,
			Email:           r.Email,
			PhoneNo:         utils.ConvertPgtypeTextToString(r.PhoneNo),
			DateOfBirth:     r.DateOfBirth.Time.Format("2006-01-02"),
			Gender:          r.Gender,
			Bio:             utils.ConvertPgtypeTextToString(r.Bio),
			ProfilePhotoUrl: utils.ConvertPgtypeTextToString(r.ProfilePhotoUrl),
			HomeCountryId:   utils.ConvertInt8ToInt64(r.HomeCountryID),
			HomeStateId:     utils.ConvertInt8ToInt64(r.HomeStateID),
			HomeCityId:      utils.ConvertInt8ToInt64(r.HomeCityID),
			IsPhoneVerified: utils.ConvertPgtypeBoolToBool(r.IsPhoneVerified),
			TrustScore:      utils.ConvertInt4ToInt64(r.TrustScore),
			Status:          utils.ConvertPgtypeTextToString(r.Status),
			PasswordHash:    r.PasswordHash,
		})
	}
	return &database_gen.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (dss *DatabaseServiceServer) CreateUser(ctx context.Context, req *database_gen.CreateUserRequest) (*database_gen.CreateUserResponse, error) {
	now := time.Now()
	user, err := dss.Queries.CreateUser(ctx, db.CreateUserParams{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PhoneNo:         utils.ConvertStringToPgtypeText(&req.PhoneNo),
		DateOfBirth:     utils.ConvertStringToPgtypeDate(&req.DateOfBirth),
		Gender:          req.Gender,
		PasswordHash:    req.PasswordHash,
		IsPhoneVerified: utils.ConvertBoolToPgtypeBool(false),
		InsertTs:        utils.ConvertTimeToPgtypeTimestamptz(now),
		ModifiedTs:      utils.ConvertTimeToPgtypeTimestamptz(now),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user in database: %v", err)
	}
	return &database_gen.CreateUserResponse{
		User: &database_gen.User{
			Id:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			PhoneNo:         utils.ConvertPgtypeTextToString(user.PhoneNo),
			DateOfBirth:     utils.ConvertPgtypeDateToString(user.DateOfBirth),
			Gender:          user.Gender,
			IsPhoneVerified: utils.ConvertPgtypeBoolToBool(user.IsPhoneVerified),
			TrustScore:      utils.ConvertInt4ToInt64(user.TrustScore),
			Status:          utils.ConvertPgtypeTextToString(user.Status),
			Bio:             utils.ConvertPgtypeTextToString(user.Bio),
			ProfilePhotoUrl: utils.ConvertPgtypeTextToString(user.ProfilePhotoUrl),
			HomeCountryId:   utils.ConvertInt8ToInt64(user.HomeCountryID),
			HomeStateId:     utils.ConvertInt8ToInt64(user.HomeStateID),
			HomeCityId:      utils.ConvertInt8ToInt64(user.HomeCityID),
			InsertTs:        utils.ConvertPgtypeTimestamptzToTimestamp(user.InsertTs),
			ModifiedTs:      utils.ConvertPgtypeTimestamptzToTimestamp(user.ModifiedTs),
		},
	}, nil
}
