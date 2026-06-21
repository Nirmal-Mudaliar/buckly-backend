package utils

import (
	"buckly-ms/gateway/api/handlers/bucket-list/dto"
	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
)

func MapBucketListItem(item *bucket_list_gen.BucketListItem) dto.BucketListItem {
	return dto.BucketListItem{
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
		InsertTs:           item.InsertTs.AsTime().String(),
		ModifiedTs:         item.ModifiedTs.AsTime().String(),
	}
}
