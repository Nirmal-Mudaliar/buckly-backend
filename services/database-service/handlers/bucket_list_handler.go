package handlers

import (
	database_gen "buckly-ms/proto/database-gen"
	db "buckly-ms/services/database-service/db/generated"
	"buckly-ms/services/database-service/utils"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (dss *DatabaseServiceServer) CreateBucketListItem(ctx context.Context, req *database_gen.CreateBucketListItemRequest) (*database_gen.CreateBucketListItemResponse, error) {
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
