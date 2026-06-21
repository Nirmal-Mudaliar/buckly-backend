package dto

type BucketListItem struct {
	Id                 int64  `json:"id"`
	UserId             int64  `json:"user_id"`
	ActivityTagId      int64  `json:"activity_tag_id"`
	CountryId          int64  `json:"country_id"`
	StateId            int64  `json:"state_id"`
	CityId             int64  `json:"city_id"`
	TimeframeStartDate string `json:"timeframe_start_date"`
	TimeframeEndDate   string `json:"timeframe_end_date"`
	Note               string `json:"note"`
	IsPublic           bool   `json:"is_public"`
	Status             string `json:"status"`
	InsertTs           string `json:"insert_ts"`
	ModifiedTs         string `json:"modified_ts"`
}

type CreateBucketListRequest struct {
	ActivityTagId      int64  `json:"activity_tag_id"`
	CountryId          int64  `json:"country_id"`
	StateId            int64  `json:"state_id"`
	CityId             int64  `json:"city_id"`
	TimeframeStartDate string `json:"timeframe_start_date"`
	TimeframeEndDate   string `json:"timeframe_end_date"`
	Note               string `json:"note"`
	IsPublic           bool   `json:"is_public"`
}

type CreateBucketListResponse struct {
	BucketListItem BucketListItem `json:"bucket_list_item"`
}

type GetBucketListItemsByUserIdResponse struct {
	BucketListItems []BucketListItem `json:"bucket_list_items"`
}

type GetBucketListItemByIdResponse struct {
	BucketListItem BucketListItem `json:"bucket_list_item"`
}

type UpdateBucketListItemRequest struct {
	Id                 int64  `json:"id"`
	ActivityTagId      int64  `json:"activity_tag_id"`
	CountryId          int64  `json:"country_id"`
	StateId            int64  `json:"state_id"`
	CityId             int64  `json:"city_id"`
	TimeframeStartDate string `json:"timeframe_start_date"`
	TimeframeEndDate   string `json:"timeframe_end_date"`
	Note               string `json:"note"`
	IsPublic           bool   `json:"is_public"`
}

type UpdateBucketListItemResponse struct {
	BucketListItem BucketListItem `json:"bucket_list_item"`
}

type DeleteBucketListItemResponse struct {
	DeletedCount int64 `json:"deleted_count"`
}
