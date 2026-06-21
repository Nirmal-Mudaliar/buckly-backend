package bucket_list

import (
	"buckly-ms/core/utils"
	"buckly-ms/gateway/api/handlers/bucket-list/dto"
	bucket_list_utils "buckly-ms/gateway/api/handlers/bucket-list/utils"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	"buckly-ms/gateway/middleware"
	"buckly-ms/gateway/models"
	"strconv"

	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	secured.GET("/items/:id", blh.GetBucketListItemById)
	secured.PUT("/update", blh.UpdateBucketListItem)
	secured.DELETE("/items/:id", blh.DeleteBucketListItem)
	secured.POST("/items/matches", blh.FindMatchesForBucketListItem)
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

// Get Bucket List Item By Id godoc
// @Summary      Get bucket list item by ID
// @Description  Fetch a bucket list item belonging to the authenticated user by its ID
// @Tags         bucket-list
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  query  int  true  "Bucket List Item ID"
// @Success      200  {object}  models.ApiResponse{data=dto.GetBucketListItemByIdResponse}
// @Router       /api/v1/bucket-list/items/:id [get]
func (blh *BucketListHandler) GetBucketListItemById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	logger := utils.GetLoggerFromContext(ctx)

	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error("Invalid Claims: ", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Invalid Claims",
		})
		return
	}

	idQuery := c.Query("id")
	id, err := strconv.ParseInt(idQuery, 10, 64)

	if err != nil {
		logger.Error("Error occurred while parsing id: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Error occurred while parsing id",
		})
		return
	}

	resp, err := blh.bucketListServiceClient.GetBucketListItemById(
		ctx,
		&bucket_list_gen.GetBucketListItemByIdRequest{
			Id:     id,
			UserId: claims.UserId,
		},
	)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, models.ApiResponse{
					Success: false,
					Message: st.Message(),
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, models.ApiResponse{
					Success: false,
					Message: "Error occured while fetching bucket list item",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occured while fetching bucket list item",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: dto.GetBucketListItemByIdResponse{
			BucketListItem: bucket_list_utils.MapBucketListItem(resp.BucketListItem),
		},
	})
}

// Update Bucket List Item godoc
// @Summary      Update bucket list item
// @Description  Update an existing bucket list item belonging to the authenticated user
// @Tags         bucket-list
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.UpdateBucketListItemRequest  true  "Bucket list item update details"
// @Success      200  {object}  models.ApiResponse{data=dto.UpdateBucketListItemResponse}
// @Router       /api/v1/bucket-list/update [put]
func (blh *BucketListHandler) UpdateBucketListItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	logger := utils.GetLoggerFromContext(ctx)

	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error("Invalid Claims: ", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Invalid Claims",
		})
		return
	}

	var payload dto.UpdateBucketListItemRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("Invalid Payload: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Payload",
		})
		return
	}

	resp, err := blh.bucketListServiceClient.UpdateBucketListItem(
		ctx,
		&bucket_list_gen.UpdateBucketListItemRequest{
			Id:                 payload.Id,
			UserId:             claims.UserId,
			ActivityTagId:      payload.ActivityTagId,
			CountryId:          payload.CountryId,
			StateId:            payload.StateId,
			CityId:             payload.CityId,
			TimeframeStartDate: payload.TimeframeStartDate,
			TimeframeEndDate:   payload.TimeframeEndDate,
			Note:               payload.Note,
			IsPublic:           payload.IsPublic,
		},
	)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, models.ApiResponse{
					Success: false,
					Message: st.Message(),
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, models.ApiResponse{
					Success: false,
					Message: "Error occured while updating bucket list item",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occured while updating bucket list item",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: dto.UpdateBucketListItemResponse{
			BucketListItem: bucket_list_utils.MapBucketListItem(resp.BucketListItem),
		},
	})
}

// DeleteBucketListItem godoc
// @Summary Delete a bucket list item
// @Description Delete a bucket list item belonging to the authenticated user
// @Tags bucket-list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Bucket List Item ID"
// @Success 200 {object} models.ApiResponse{data=dto.DeleteBucketListItemResponse}
// @Router /api/v1/bucket-list/items/:id [delete]
func (blh *BucketListHandler) DeleteBucketListItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	logger := utils.GetLoggerFromContext(ctx)

	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error("Invalid Claims: ", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Invalid Claims",
		})
		return
	}

	idQuery := c.Query("id")
	id, err := strconv.ParseInt(idQuery, 10, 64)
	if err != nil {
		logger.Error("Error occured while parsing id from query params: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Query Params",
		})
		return
	}

	resp, err := blh.bucketListServiceClient.DeleteBucketListItem(
		ctx,
		&bucket_list_gen.DeleteBucketListItemRequest{
			Id:     id,
			UserId: claims.UserId,
		},
	)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, models.ApiResponse{
					Success: false,
					Message: "Error occurred while deleting bucket list",
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, models.ApiResponse{
					Success: false,
					Message: "Error occurred while deleting bucket list",
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: dto.DeleteBucketListItemResponse{
			DeletedCount: resp.DeletedCount,
		},
	})
}

// FindMatchesForBucketListItem godoc
// @Summary Find matching bucket list items
// @Description Find bucket list items from other users that match the specified activity, city, and timeframe
// @Tags bucket-list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.FindMatchesForBucketListItemRequest true "Match criteria"
// @Success 200 {object} models.ApiResponse{data=dto.FindMatchesForBucketListItemResponse}
// @Router /api/v1/bucket-list/items/matches [post]
func (blh *BucketListHandler) FindMatchesForBucketListItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	logger := utils.GetLoggerFromContext(ctx)

	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error("Invalid Claims: ", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Invalid Claims",
		})
		return
	}

	var payload dto.FindMatchesForBucketListItemRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("Invalid Payload:", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid Payload",
		})
		return
	}

	resp, err := blh.bucketListServiceClient.FindMatchesForBucketListItem(
		ctx,
		&bucket_list_gen.FindMatchesForBucketListItemRequest{
			ActivityTagId:      payload.ActivityTagId,
			CityId:             payload.CityId,
			TimeframeStartDate: payload.TimeframeStartDate,
			TimeframeEndDate:   payload.TimeframeEndDate,
			UserId:             claims.UserId,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Error occurred while finding matches for bucket list item",
		})
		return
	}

	bucketListItems := make([]dto.BucketListItem, 0, len(resp.BucketListItems))
	for _, item := range resp.BucketListItems {
		bucketListItems = append(bucketListItems, bucket_list_utils.MapBucketListItem(item))
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: dto.FindMatchesForBucketListItemResponse{
			BucketListItems: bucketListItems,
		},
	})
}
