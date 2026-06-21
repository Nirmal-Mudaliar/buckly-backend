package handlers

import (
	core_utils "buckly-ms/core/utils"
	database_gen "buckly-ms/proto/database-gen"
	db "buckly-ms/services/database-service/db/generated"
	"buckly-ms/services/database-service/utils"
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (dss *DatabaseServiceServer) CreateBucketListItem(ctx context.Context, req *database_gen.CreateBucketListItemRequest) (*database_gen.CreateBucketListItemResponse, error) {
	logger := core_utils.GetLoggerFromContext(ctx)
	now := time.Now()
	item, err := dss.Queries.CreateBucketListItem(ctx, db.CreateBucketListItemParams{
		UserID:             req.UserId,
		ActivityTagID:      req.ActivityTagId,
		CountryID:          req.CountryId,
		StateID:            req.StateId,
		CityID:             req.CityId,
		TimeframeStartDate: utils.ConvertStringToPgtypeDate(&req.TimeframeStartDate),
		TimeframeEndDate:   utils.ConvertStringToPgtypeDate(&req.TimeframeEndDate),
		Note:               utils.ConvertStringToPgtypeText(&req.Note),
		IsPublic:           req.IsPublic,
		Status:             "planned",
		InsertTs:           utils.ConvertTimeToPgtypeTimestamptz(now),
		ModifiedTs:         utils.ConvertTimeToPgtypeTimestamptz(now),
	})

	if err != nil {
		logger.Error("Error creating bucket list item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create bucket list item: %v", err)
	}

	return &database_gen.CreateBucketListItemResponse{
		BucketListItem: mapBucketListItem(item),
	}, nil
}

func mapBucketListItem(item db.BucketListItem) *database_gen.BucketListItem {
	return &database_gen.BucketListItem{
		Id:                 item.ID,
		UserId:             item.UserID,
		ActivityTagId:      item.ActivityTagID,
		CountryId:          item.CountryID,
		StateId:            item.StateID,
		CityId:             item.CityID,
		TimeframeStartDate: utils.ConvertPgtypeDateToString(item.TimeframeStartDate),
		TimeframeEndDate:   utils.ConvertPgtypeDateToString(item.TimeframeEndDate),
		Note:               utils.ConvertPgtypeTextToString(item.Note),
		IsPublic:           item.IsPublic,
		Status:             item.Status,
		InsertTs:           utils.ConvertPgtypeTimestamptzToTimestamp(item.InsertTs),
		ModifiedTs:         utils.ConvertPgtypeTimestamptzToTimestamp(item.ModifiedTs),
	}
}

func (dss *DatabaseServiceServer) GetBucketListItemsByUserId(
	ctx context.Context,
	req *database_gen.GetBucketListItemsByUserIdRequest,
) (
	*database_gen.GetBucketListItemsByUserIdResponse,
	error,
) {
	logger := core_utils.GetLoggerFromContext(ctx)
	resp, err := dss.Queries.GetBucketListItemsByUserId(
		ctx,
		req.UserId,
	)

	if err != nil {
		logger.Error("Error occurred while fetching bucket list items", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to fetch bucket list items: %v", err)
	}

	bucketListItems := make([]*database_gen.BucketListItem, len(resp))
	for i, item := range resp {
		bucketListItems[i] = mapBucketListItem(item)
	}

	return &database_gen.GetBucketListItemsByUserIdResponse{
		BucketListItems: bucketListItems,
	}, nil
}
