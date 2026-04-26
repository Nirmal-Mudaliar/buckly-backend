package handlers

import (
	core_utils "buckly-ms/core/utils"
	database_gen "buckly-ms/proto/database-gen"
	"buckly-ms/services/database-service/utils"
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (dss *DatabaseServiceServer) GetAllUsers(ctx context.Context, empty *emptypb.Empty) (*database_gen.GetAllUsersResponse, error) {
	logger := core_utils.GetLoggerFromContext(ctx)
	rows, err := dss.Queries.GetAllUsers(ctx)

	if err != nil {
		logger.Error("Error fetching users from database", zap.Error(err))
		return nil, err
	}
	users := make([]*database_gen.Users, 0, len(rows))
	for _, r := range rows {
		users = append(users, &database_gen.Users{
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
		})
	}
	return &database_gen.GetAllUsersResponse{
		Users: users,
	}, nil
}
