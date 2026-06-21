package bucket_list

import (
	"buckly-ms/core/utils"
	"buckly-ms/gateway/api/handlers/bucket-list/dto"
	bucket_list_utils "buckly-ms/gateway/api/handlers/bucket-list/utils"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	"buckly-ms/gateway/middleware"
	"buckly-ms/gateway/models"

	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BucketListHandler struct {
	config                  *config.GatewayConfig
	bucketListServiceClient bucket_list_gen.BucketListServiceClient
}

func NewBucketListHandler(
	config *config.GatewayConfig,
	bucketListServiceClient bucket_list_gen.BucketListServiceClient,
) contracts.RouteRegistrar {
	return &BucketListHandler{
		config:                  config,
		bucketListServiceClient: bucketListServiceClient,
	}
}

func (blh *BucketListHandler) RegisterRoutes(r *gin.Engine) {
	secured := r.Group("/api/v1/bucket-list", middleware.TokenMiddleware(blh.config.JWTSecret))
	secured.POST("/create", blh.CreateBucketList)
	secured.GET("/items", blh.GetBucketListItemsByUserId)
}

// Create Bucket List godoc
// @Summary      Create a new bucket list item
// @Description  Creates a new bucket list item with activity, location, and timeframe details for the authenticated user
// @Tags         bucket-list
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateBucketListRequest  true  "Bucket list item details"
// @Success      200  {object}  models.ApiResponse{data=dto.CreateBucketListResponse}
// @Router       /api/v1/bucket-list/create [post]
func (blh *BucketListHandler) CreateBucketList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	claims, err := utils.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Invalid Claims",
		})
		return
	}

	var payload dto.CreateBucketListRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Payload",
		})
		return
	}

	resp, err := blh.bucketListServiceClient.CreateBucketListItem(ctx, &bucket_list_gen.CreateBucketListItemRequest{
		UserId:             claims.UserId,
		ActivityTagId:      payload.ActivityTagId,
		CountryId:          payload.CountryId,
		StateId:            payload.StateId,
		CityId:             payload.CityId,
		TimeframeStartDate: payload.TimeframeStartDate,
		TimeframeEndDate:   payload.TimeframeEndDate,
		Note:               payload.Note,
		IsPublic:           payload.IsPublic,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Failed to create bucket list item",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: dto.CreateBucketListResponse{
			BucketListItem: bucket_list_utils.MapBucketListItem(resp.BucketListItem),
		},
	})

}

// Get Bucket List Items By User Id godoc
// @Summary      Get bucket list items for the authenticated user
// @Description  Fetch all bucket list items belonging to the authenticated user
// @Tags         bucket-list
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.ApiResponse{data=dto.CreateBucketListResponse}
// @Router       /api/v1/bucket-list/items [get]
func (blh *BucketListHandler) GetBucketListItemsByUserId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	logger := utils.GetLoggerFromContext(ctx)

	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error("Invalid Claims", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Invalid Claims",
		})
		return
	}

	resp, err := blh.bucketListServiceClient.GetBucketListItemsByUserId(
		ctx,
		&bucket_list_gen.GetBucketListItemsByUserIdRequest{
			UserId: claims.UserId,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occurred while fetching bucket list items",
		})
		return
	}

	bucketListItems := make([]dto.BucketListItem, 0, len(resp.BucketListItems))
	for _, item := range resp.BucketListItems {
		bucketListItems = append(bucketListItems, bucket_list_utils.MapBucketListItem(item))
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: dto.GetBucketListItemsByUserIdResponse{
			BucketListItems: bucketListItems,
		},
	})
}
