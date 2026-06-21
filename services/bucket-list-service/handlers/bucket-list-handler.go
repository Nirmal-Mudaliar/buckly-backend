package handlers

import (
	"buckly-ms/core/utils"
	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	database_gen "buckly-ms/proto/database-gen"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (blss *BucketListServiceServer) CreateBucketListItem(ctx context.Context, req *bucket_list_gen.CreateBucketListItemRequest) (*bucket_list_gen.CreateBucketListItemResponse, error) {
	logger := utils.GetLoggerFromContext(ctx)

	resp, err := blss.DatabaseServiceClient.CreateBucketListItem(ctx, &database_gen.CreateBucketListItemRequest{
		UserId:             req.UserId,
		ActivityTagId:      req.ActivityTagId,
		CountryId:          req.CountryId,
		StateId:            req.StateId,
		CityId:             req.CityId,
		TimeframeStartDate: req.TimeframeStartDate,
		TimeframeEndDate:   req.TimeframeEndDate,
		Note:               req.Note,
		IsPublic:           req.IsPublic,
	})
	if err != nil {
		logger.Error("Error creating bucket list item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create bucket list item: %v", err)
	}

	return &bucket_list_gen.CreateBucketListItemResponse{
		BucketListItem: mapBucketListItem(resp.BucketListItem),
	}, nil
}

func mapBucketListItem(item *database_gen.BucketListItem) *bucket_list_gen.BucketListItem {
	if item == nil {
		return nil
	}
	return &bucket_list_gen.BucketListItem{
		Id:                 item.Id,
		UserId:             item.UserId,
		ActivityTagId:      item.ActivityTagId,
		CountryId:          item.CountryId,
		StateId:            item.StateId,
		CityId:             item.CityId,
		TimeframeStartDate: item.TimeframeStartDate,
		TimeframeEndDate:   item.TimeframeEndDate,
		Note:               item.Note,
		IsPublic:           item.IsPublic,
		Status:             item.Status,
		InsertTs:           item.InsertTs,
		ModifiedTs:         item.ModifiedTs,
	}
}
