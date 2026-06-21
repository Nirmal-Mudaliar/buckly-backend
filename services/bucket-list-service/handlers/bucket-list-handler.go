package handlers

import (
	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	database_gen "buckly-ms/proto/database-gen"
	"context"
)

func (blss *BucketListServiceServer) CreateBucketListItem(ctx context.Context, req *bucket_list_gen.CreateBucketListItemRequest) (*bucket_list_gen.CreateBucketListItemResponse, error) {
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
		return nil, err
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

func (blss *BucketListServiceServer) GetBucketListItemsByUserId(
	ctx context.Context,
	req *bucket_list_gen.GetBucketListItemsByUserIdRequest,
) (
	*bucket_list_gen.GetBucketListItemsByUserIdResponse,
	error,
) {
	resp, err := blss.DatabaseServiceClient.GetBucketListItemsByUserId(
		ctx,
		&database_gen.GetBucketListItemsByUserIdRequest{
			UserId: req.UserId,
		},
	)

	if err != nil {
		return nil, err
	}

	bucketListItems := make([]*bucket_list_gen.BucketListItem, 0, len(resp.BucketListItems))
	for _, item := range resp.BucketListItems {
		bucketListItems = append(bucketListItems, mapBucketListItem(item))
	}

	return &bucket_list_gen.GetBucketListItemsByUserIdResponse{
		BucketListItems: bucketListItems,
	}, nil
}

func (blss *BucketListServiceServer) GetBucketListItemById(
	ctx context.Context,
	req *bucket_list_gen.GetBucketListItemByIdRequest,
) (
	*bucket_list_gen.GetBucketListItemByIdResponse,
	error,
) {
	resp, err := blss.DatabaseServiceClient.GetBucketListItemById(
		ctx,
		&database_gen.GetBucketListItemByIdRequest{
			Id:     req.Id,
			UserId: req.UserId,
		},
	)

	if err != nil {
		return nil, err
	}

	return &bucket_list_gen.GetBucketListItemByIdResponse{
		BucketListItem: mapBucketListItem(resp.BucketListItem),
	}, nil
}

func (blss *BucketListServiceServer) UpdateBucketListItem(
	ctx context.Context,
	req *bucket_list_gen.UpdateBucketListItemRequest,
) (
	*bucket_list_gen.UpdateBucketListItemResponse,
	error,
) {
	resp, err := blss.DatabaseServiceClient.UpdateBucketListItem(
		ctx,
		&database_gen.UpdateBucketListItemRequest{
			Id:                 req.Id,
			UserId:             req.UserId,
			ActivityTagId:      req.ActivityTagId,
			CountryId:          req.CountryId,
			StateId:            req.StateId,
			CityId:             req.CityId,
			TimeframeStartDate: req.TimeframeStartDate,
			TimeframeEndDate:   req.TimeframeEndDate,
			Note:               req.Note,
			IsPublic:           req.IsPublic,
		},
	)
	if err != nil {
		return nil, err
	}

	return &bucket_list_gen.UpdateBucketListItemResponse{
		BucketListItem: mapBucketListItem(resp.BucketListItem),
	}, nil
}

func (blss *BucketListServiceServer) DeleteBucketListItem(
	ctx context.Context,
	req *bucket_list_gen.DeleteBucketListItemRequest,
) (
	*bucket_list_gen.DeleteBucketListItemResponse,
	error,
) {
	resp, err := blss.DatabaseServiceClient.DeleteBucketListItem(
		ctx,
		&database_gen.DeleteBucketListItemRequest{
			Id:     req.Id,
			UserId: req.UserId,
		},
	)

	if err != nil {
		return nil, err
	}

	return &bucket_list_gen.DeleteBucketListItemResponse{
		DeletedCount: resp.DeletedCount,
	}, nil
}
